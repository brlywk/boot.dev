package commands

import (
	"brlywk/bootdev/pokedex/pokeapi"
	"errors"
	"fmt"
)

// Explore a single area and print all pokemon encounters in that area
func CommandExplore(config *pokeapi.ApiConfig, areaName string) error {
	if areaName == "" {
		return errors.New("Please provide an area to explore")
	}

	resp, err := pokeapi.FetchExplorationData(config.LocationAreaUrl, areaName, config.Cache)
	if err != nil {
		return err
	}

	for _, encounter := range resp.Encounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return nil
}
