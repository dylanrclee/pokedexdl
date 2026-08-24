package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Print("Hello, World!")
}

func cleanInput(text string) []string {
	lowertext := strings.ToLower(text)
	return strings.Fields(lowertext)
}
