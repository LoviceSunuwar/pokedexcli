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
	callback    func() error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
}

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
				cmd.callback()
			} else {
				fmt.Println("Unknown command")
			}
		} else {
			fmt.Println("It's empty you mortal")
		}

	}

}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	spaceRemoved := strings.TrimSpace(lowerText)
	return strings.Fields(spaceRemoved)
}
