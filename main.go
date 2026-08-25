package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		scan_text := scanner.Text()
		ct := cleanInput(scan_text)
		fmt.Printf("Your command was: %s\n", ct[0])
	}
}

func cleanInput(text string) []string {
	lowertext := strings.ToLower(text)
	return strings.Fields(lowertext)
}
