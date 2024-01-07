package commands

import (
	"brlywk/bootdev/pokedex/pokeapi"
	"fmt"
)

// Displays the names of 20 location areas in the Pokemon world
//
// Subsequent calls should displays the next 20 locations etc.
func CommandMap(config *pokeapi.ApiConfig, _ string) error {
	// otherwise fetch
	resp, err := pokeapi.FetchLocationData(config.NextLocation, config.Cache)
	if err != nil {
		return err
	}

	config.NextLocation = resp.Next
	config.PreviousLocation = resp.Previous

	for _, location := range resp.Results {
		fmt.Printf("\t %v\n", location.Name)
	}

	return nil
}

// Displays the previous 20 locations
//
// Prints an error if there are no more locations to go back to
func CommandMapB(config *pokeapi.ApiConfig, _ string) error {
	if config.PreviousLocation == "" {
		return fmt.Errorf("You are on the first page")
	}

	resp, err := pokeapi.FetchLocationData(config.PreviousLocation, config.Cache)
	if err != nil {
		return err
	}

	config.NextLocation = resp.Next
	config.PreviousLocation = resp.Previous

	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}

	return nil
}
