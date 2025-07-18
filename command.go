package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type config struct {
	next     string
	previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

const baseUrl string = "https://pokeapi.co/api/v2/"

const locationArea string = "location-area/"

var configs config

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
		"map": {
			name:        "map",
			description: "Displays the list of map",
			callback:    fetchMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Go back to the previous map",
			callback:    fetchmapB,
		},
	}
}

func commandExit(*config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(*config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func fetchMap(cfg *config) error {
	var fullUrl string

	if cfg.next == "" {
		fullUrl = baseUrl + locationArea
	} else {
		fullUrl = cfg.next
	}

	res, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Error fetching the location")
		log.Fatal(err)
	}

	bytBody, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s ", res.StatusCode, bytBody)
	}
	if err != nil {
		log.Fatal(err)
	}

	resultLocationArea := locationAreaS{}

	err = json.Unmarshal(bytBody, &resultLocationArea)
	if err != nil {
		log.Fatalf("Unmarshal failed %v", err)
	}

	cfg.next = resultLocationArea.Next
	cfg.previous = resultLocationArea.Previous

	locationResults := resultLocationArea.Results

	for _, x := range locationResults {
		fmt.Println(x.Name)
	}

	return nil
}

func fetchmapB(cfg *config) error {
	var fullUrl string

	if cfg.previous == "" {
		fmt.Print("You are on the first page")
		return nil
	} else {
		fullUrl = cfg.previous
	}

	res, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Error fetching the location")
		log.Fatal(err)
	}

	bytBody, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s ", res.StatusCode, bytBody)
	}
	if err != nil {
		log.Fatal(err)
	}

	resultLocationArea := locationAreaS{}

	err = json.Unmarshal(bytBody, &resultLocationArea)
	if err != nil {
		log.Fatalf("Unmarshal failed %v", err)
	}

	cfg.next = resultLocationArea.Next
	cfg.previous = resultLocationArea.Previous

	locationResults := resultLocationArea.Results

	for _, x := range locationResults {
		fmt.Println(x.Name)
	}

	return nil
}
