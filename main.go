package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dylanrclee/pokedexdl/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	client       internal.Client
	command_list map[string]cliCommand
	Next         string
	Previous     string
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
	}

	comm_reg := &config{}
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
			err := val.callback(comms)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Print("Unknown command\n\n")
		}
	}
}

func commandExit(cur_config *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n\n")
	os.Exit(0)
	return nil
}

func commandHelp(cur_config *config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n")
	fmt.Print("help: Displays a help message\n")
	fmt.Print("map: Displays the first or next 20 locations\n")
	fmt.Print("mapb: Displays the previous 20 locations\n")
	fmt.Print("exit: Exit the Pokedex\n\n")
	return nil
}

func commandMap(cur_config *config) error {
	var call_url string
	if cur_config.Next == "" {
		call_url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		call_url = cur_config.Next
	}
	var city_struct internal.Pokestruct
	err := internal.Pokeapi(call_url, &city_struct, cur_config.client.Cache)
	if err != nil {
		return err
	}
	print("\n")
	for _, val := range city_struct.Results {
		fmt.Println(val.Name)
	}
	fmt.Print("\n\n")
	cur_config.Previous = city_struct.Previous
	cur_config.Next = city_struct.Next
	return nil
}

func commandMapb(cur_config *config) error {
	var call_url string
	if cur_config.Previous == "" || cur_config.Previous == "null" {
		fmt.Print("you're on the first page\n\n")
		return nil
	} else {
		call_url = cur_config.Previous
	}
	var city_struct internal.Pokestruct
	internal.Pokeapi(call_url, &city_struct, cur_config.client.Cache)
	for _, val := range city_struct.Results {
		fmt.Println(val.Name)
	}
	fmt.Print("\n\n")
	cur_config.Previous = city_struct.Previous
	cur_config.Next = city_struct.Next
	return nil
}

func cleanInput(text string) []string {
	lowertext := strings.ToLower(text)
	return strings.Fields(lowertext)
}
