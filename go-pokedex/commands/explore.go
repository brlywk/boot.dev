package commands

import (
	"brlywk/bootdev/pokedex/pokeapi"
	"errors"
	"fmt"
)

func CommandExplore(config *pokeapi.ApiConfig, areaName string) error {
	if areaName == "" {
		return errors.New("Please provide an area to explore")
	}

	explorationApi := fmt.Sprintf("%v/%v", config.BaseUrl, config.LocationPath)
	resp, err := pokeapi.FetchExplorationData(explorationApi, areaName, config.Cache)
	if err != nil {
		return err
	}

	for _, encounter := range resp.Encounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return nil
}
