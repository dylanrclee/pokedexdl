package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	type cliCommand struct {
		name        string
		description string
		callback    func() error
	}

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

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		scan_text := scanner.Text()
		ct := cleanInput(scan_text)
		val, ok := commands[ct[0]]
		if ok {
			val.callback()
		} else {
			fmt.Print("Unknown command\n")
		}
	}
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return fmt.Errorf("No error, exit successful")
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n")
	fmt.Print("help: Displays a help message\n")
	fmt.Print("exit: Exit the Pokedex\n")
	return fmt.Errorf("No error, help successful")
}

func cleanInput(text string) []string {
	lowertext := strings.ToLower(text)
	return strings.Fields(lowertext)
}
