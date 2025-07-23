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
				if cmdName == "explore" {
					if len(words) < 2 {
						fmt.Println("Usage: explore <area-name>")
						continue
					}
					configs.areaName = words[1]
					cmd.callback(&configs, words[1:])
				} else {
					cmd.callback(&configs, words[1:])
				}
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
