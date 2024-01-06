package pokeapi

import "brlywk/bootdev/pokedex/cache"

// ----- Config ----------------------------------

// Hold info about where to start fetching
type ApiConfig struct {
	LocationAreaUrl  string
	NextLocation     string
	PreviousLocation string
	Cache            *cache.Cache
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
