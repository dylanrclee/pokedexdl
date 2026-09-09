package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dylanrclee/pokedexdl/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	client         internal.Client
	caught_pokemon map[string]internal.Pokemon
	command_list   map[string]cliCommand
	Next           string
	Previous       string
}

func main() {

	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},

		"help": {
			name:        "help",
			description: "explains pokedex and what commands can be used",
			callback:    commandHelp,
		},

		"map": {
			name:        "map",
			description: "displays next 20 locations of the pokemon world",
			callback:    commandMap,
		},

		"mapb": {
			name:        "mapb",
			description: "displays previous 20 location of the pokemon world",
			callback:    commandMapb,
		},

		"explore": {
			name:        "explore",
			description: "displays the pokemon that can be encountered at the inputed location",
			callback:    commandexplore,
		},

		"catch": {
			name:        "catch",
			description: "shows catching messages depending on if pokemon was caught based on its catch rate",
			callback:    commandcatch,
		},
	}

	comm_reg := &config{}
	comm_reg.caught_pokemon = make(map[string]internal.Pokemon)
	comm_reg.command_list = commands
	comm_reg.client.Cache = internal.NewCache(5 * time.Second)
	REPLloop((comm_reg))
}

func REPLloop(comms *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		scan_text := scanner.Text()
		ct := cleanInput(scan_text)
		if len(ct) == 0 {
			fmt.Print("No command entered\n\n")
			continue
		}
		val, ok := comms.command_list[ct[0]]
		if ok {
			err := val.callback(comms, ct)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Print("Unknown command\n\n")
		}
	}
}

func commandExit(cur_config *config, _ []string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n\n")
	os.Exit(0)
	return nil
}

func commandHelp(cur_config *config, _ []string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n")
	fmt.Print("help: Displays a help message\n")
	fmt.Print("map: Displays the first or next 20 locations\n")
	fmt.Print("mapb: Displays the previous 20 locations\n")
	fmt.Print("explore *location*: Displays the pokemon found at the inputed location\n")
	fmt.Print("exit: Exit the Pokedex\n\n")
	return nil
}

func commandMap(cur_config *config, _ []string) error {
	var call_url string
	if cur_config.Next == "" {
		call_url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		call_url = cur_config.Next
	}
	var city_struct internal.Pokestruct
	err := cur_config.client.ListLocations(call_url, &city_struct, cur_config.client.Cache)
	if err != nil {
		return err
	}
	fmt.Print("\n")
	for _, val := range city_struct.Results {
		fmt.Println(val.Name)
	}
	fmt.Print("\n\n")
	cur_config.Previous = city_struct.Previous
	cur_config.Next = city_struct.Next
	return nil
}

func commandMapb(cur_config *config, _ []string) error {
	var call_url string
	if cur_config.Previous == "" || cur_config.Previous == "null" {
		fmt.Print("you're on the first page\n\n")
		return nil
	} else {
		call_url = cur_config.Previous
	}
	var city_struct internal.Pokestruct
	err := cur_config.client.ListLocations(call_url, &city_struct, cur_config.client.Cache)
	if err != nil {
		return err
	}
	fmt.Print("\n")
	for _, val := range city_struct.Results {
		fmt.Println(val.Name)
	}
	fmt.Print("\n\n")
	cur_config.Previous = city_struct.Previous
	cur_config.Next = city_struct.Next
	return nil
}

func commandexplore(cur_config *config, location []string) error {
	if len(location) < 2 {
		fmt.Print("no location entered\n\n")
		return nil
	}
	var location_pokemon internal.Cityinfo
	err := cur_config.client.GetLocationPkmn(location[1], &location_pokemon, cur_config.client.Cache)
	if err != nil {
		fmt.Print("\nPossible misspelling of location\n")
		return err
	}
	fmt.Print("\n")
	for _, val := range location_pokemon.Pkmn_enc {
		fmt.Println(val.Pokemon.Name)
	}
	fmt.Print("\n\n")
	return nil
}

func commandcatch(cur_config *config, catching_pokemon []string) error {
	if len(catching_pokemon) < 2 {
		fmt.Print("no pokemon entered\n\n")
		return nil
	}

	var pokemoninfo internal.Pokemon
	err := cur_config.client.GetPokemoninfo(catching_pokemon[1], &pokemoninfo, cur_config.client.Cache)
	if err != nil {
		fmt.Print("\nPossible misspelling of pokemon\n")
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemoninfo.Name)
	ok := rand.Intn(400) >= pokemoninfo.Base_exp

	if !ok {
		fmt.Printf("%s escaped!\n\n", pokemoninfo.Name)
		return nil
	} else {
		fmt.Printf("%s was caught!\n\n", pokemoninfo.Name)
		cur_config.caught_pokemon[pokemoninfo.Name] = pokemoninfo
	}
	return nil
}

func cleanInput(text string) []string {
	lowertext := strings.ToLower(text)
	return strings.Fields(lowertext)
}
