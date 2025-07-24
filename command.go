package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/lovicesunuwar/pokedexcli/internal/pokecache"
)

type config struct {
	next        string
	previous    string
	cache       pokecache.Cache
	areaName    string
	pokemonName string
	pokedex     map[string]PokemonPokeball
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

const baseUrl string = "https://pokeapi.co/api/v2/"

const locationArea string = "location-area/"

const catchPokemon string = "pokemon/"

var configs config

var commands map[string]cliCommand

func init() {

	configs.cache = pokecache.NewCache(5 * time.Minute)
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
		"explore": {
			name:        "explore",
			description: "Explore the location via name",
			callback:    explore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon to add on the user's Pokedex",
			callback:    catch,
		},
	}
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func catch(cfg *config, args []string) error {

	var fullUrl string
	resultPokemon := PokemonPokeball{}
	var bytBody []byte

	if len(args) < 1 {
		fmt.Println("Insufficient Commands")
		return fmt.Errorf("The pokemon name is missing...")
	} else {
		fullUrl = baseUrl + catchPokemon + args[0] + "/"
		if cachedData, found := cfg.cache.Get(fullUrl); found {
			bytBody = cachedData
		} else {
			res, err := http.Get(fullUrl)
			if err != nil {
				fmt.Println("Error Fecthing the response for Pokemon")
				log.Fatal(err)
			}
			defer res.Body.Close()
			bytBody, err = io.ReadAll(res.Body)
			if err != nil {
				fmt.Println("Error reading the body of the response of pokemon")
				log.Fatal(err)
			}
			if res.StatusCode > 299 {
				log.Fatalf("Response failed with status code: %d and\nbody: %s ", res.StatusCode, bytBody)
			}
			cfg.cache.Add(fullUrl, bytBody)
		}
	}

	err := json.Unmarshal(bytBody, &resultPokemon)
	if err != nil {
		fmt.Println("Error unmarshalling the data form the bytes")
		log.Fatal(err)
	}

	maxVal := 2 * resultPokemon.BaseExperience
	roll := rand.Intn(maxVal)
	fmt.Printf("Throwing a Pokeball at %v...\n", resultPokemon.Name)
	if roll < resultPokemon.BaseExperience {
		fmt.Printf("%v escaped!\n", resultPokemon.Name)
	} else {
		fmt.Printf("%v was caught!\n", resultPokemon.Name)
		if cfg.pokedex == nil {
			cfg.pokedex = make(map[string]PokemonPokeball)
		}
		cfg.pokedex[resultPokemon.Name] = resultPokemon
	}

	return nil

}

func explore(cfg *config, args []string) error {

	var fullUrl string
	resultexpArea := ExploreArea{}
	var bytBody []byte
	if len(args) < 1 {
		fmt.Println("The args string in explore is empty")
		return fmt.Errorf("missing location area name")
	} else {
		fullUrl = baseUrl + locationArea + args[0] + "/"
		if cachedData, found := cfg.cache.Get(fullUrl); found {
			fmt.Println("Using Cached Data")
			bytBody = cachedData
		} else {
			fmt.Println("Fetching the Data for explore")
			res, err := http.Get(fullUrl)
			if err != nil {
				fmt.Println("Error fetching the location-area details")
				log.Fatal(err)
			}

			defer res.Body.Close()

			bytBody, err = io.ReadAll(res.Body)
			if err != nil {
				fmt.Println("Error reading the response body")
				log.Fatal(err)
			}
			if res.StatusCode > 299 {
				log.Fatalf("Response failed with status code: %d and\nbody: %s ", res.StatusCode, bytBody)
			}
			cfg.cache.Add(fullUrl, bytBody)
		}
	}

	err := json.Unmarshal(bytBody, &resultexpArea)
	if err != nil {
		log.Fatalf("Unmarshal failed %v", err)
	}

	expAreaResults := resultexpArea.PokemonEncounters
	// fmt.Printf("DEBUG: bytBody length: %d\n", len(bytBody))                              // Add this line
	// fmt.Printf("DEBUG: bytBody content (first 200 chars):\n%s\n", string(bytBody[:200])) // And this line (shows start of content)
	fmt.Printf("Exploring %v \n", args[0])
	fmt.Println("Found Pokemon:")

	for _, x := range expAreaResults {

		fmt.Printf("- %v \n", x.Pokemon.Name)
	}

	return nil
}

func fetchMap(cfg *config, args []string) error {
	var fullUrl string
	var bytBody []byte

	if cfg.next == "" {
		fullUrl = baseUrl + locationArea
	} else {
		fullUrl = cfg.next
	}

	if cachedData, found := cfg.cache.Get(fullUrl); found {
		fmt.Println("Using cached Data")
		bytBody = cachedData
	} else {
		fmt.Println("Fetching from API")
		res, err := http.Get(fullUrl)
		if err != nil {
			fmt.Println("Error fetching the location")
			log.Fatal(err)
		}

		bytBody, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s ", res.StatusCode, bytBody)
		}
		if err != nil {
			log.Fatal(err)
		}
		cfg.cache.Add(fullUrl, bytBody)
	}

	resultLocationArea := locationAreaS{}

	err := json.Unmarshal(bytBody, &resultLocationArea)
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

func fetchmapB(cfg *config, args []string) error {
	var fullUrl string
	var bytBody []byte

	if cfg.previous == "" {
		fmt.Print("You are on the first page")
		return nil
	} else {
		fullUrl = cfg.previous
	}

	if cachedData, found := cfg.cache.Get(fullUrl); found {
		fmt.Println("Using Cached Data on mapB")
		bytBody = cachedData
	} else {
		res, err := http.Get(fullUrl)
		if err != nil {
			fmt.Println("Error fetching the location")
			log.Fatal(err)
		}

		bytBody, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s ", res.StatusCode, bytBody)
		}
		if err != nil {
			log.Fatal(err)
		}

		cfg.cache.Add(fullUrl, bytBody)
	}

	resultLocationArea := locationAreaS{}

	err := json.Unmarshal(bytBody, &resultLocationArea)
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
