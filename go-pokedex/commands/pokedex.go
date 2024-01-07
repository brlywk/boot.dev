package commands

import (
	"brlywk/bootdev/pokedex/pokeapi"
	"errors"
	"fmt"
)

func CommandPokedex(config *pokeapi.ApiConfig, _ string) error {
	caught := *config.CaughtPokemon

	if len(caught) == 0 {
		return errors.New("You have not caught any Pokemon yet. Go catch them all!")
	}

	for _, pokemon := range caught {
		fmt.Printf("\t󰐝 %v\n", pokemon.Name)
	}

	return nil
}
