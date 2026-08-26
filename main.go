package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	command_list map[string]cliCommand
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
	}

	comm_reg := &config{commands}
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
			val.callback(comms)
		} else {
			fmt.Print("Unknown command\n\n")
		}
	}
}

func commandExit(*config) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n\n")
	os.Exit(0)
	return nil
}

func commandHelp(*config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n")
	fmt.Print("help: Displays a help message\n")
	fmt.Print("exit: Exit the Pokedex\n\n")
	return nil
}

func cleanInput(text string) []string {
	lowertext := strings.ToLower(text)
	return strings.Fields(lowertext)
}
