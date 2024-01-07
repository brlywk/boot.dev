package pokeapi

import (
	"brlywk/bootdev/pokedex/cache"
)

// ----- Config ----------------------------------

// Hold info about where to start fetching
type ApiConfig struct {
	LocationAreaUrl  string
	PokemonUrl       string
	NextLocation     string
	PreviousLocation string
	Cache            *cache.Cache
	CaughtPokemon    *map[string]PokemonResponse
}

// ----- Locations -------------------------------

// Holds a single Location area
type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// Holds the full location-area API response
type LocationResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

// ----- Explore Location ------------------------

type ExplorationResponse struct {
	Id         int          `json:"id"`
	Location   LocationArea `json:"location"`
	Encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// ----- Catch Pokemon ---------------------------

type PokemonResponse struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	}
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
		} `json:"ability"`
	} `json:"abilities"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}
