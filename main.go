package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// fmt.Println("Hello, World!")
	NewScanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		NewScanner.Scan()
		newInput := NewScanner.Text()
		lowerText := strings.ToLower(newInput)
		words := strings.Fields(lowerText)
		if len(words) > 0 {
			cmdName := words[0]
			cmd, exists := commands[cmdName]
			if exists {
				cmd.callback(&configs)
			} else {
				fmt.Println("Unknown command")
			}
		} else {
			fmt.Println("It's empty you mortal")
		}
	}

}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	spaceRemoved := strings.TrimSpace(lowerText)
	return strings.Fields(spaceRemoved)
}
