package main

type ExploreArea struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}

type PokemonEncounters struct {
	Pokemon PEPokemons `json:"pokemon"`
}

type PEPokemons struct {
	Name string `json:"name"`
}
