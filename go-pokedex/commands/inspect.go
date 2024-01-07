package commands

import (
	"brlywk/bootdev/pokedex/pokeapi"
	"brlywk/bootdev/pokedex/utils"
	"errors"
	"fmt"
)

func CommandInspect(config *pokeapi.ApiConfig, pokemon string) error {
	if pokemon == "" {
		return errors.New("Please provide an area to explore")
	}

	caught := *config.CaughtPokemon
	pokemonResp, found := caught[pokemon]

	if !found {
		fmt.Printf("\tYou have not caught a %v%v%v yet.\n", utils.PromptColorOrange, pokemon, utils.PromptColorReset)
		return nil
	}

	// Print the pokemon...
	fmt.Printf("\tName:\t%v%v%v\n", utils.PromptColorOrange, pokemonResp.Name, utils.PromptColorReset)
	fmt.Printf("\tHeight:\t%v\n\tWeight:\t%v\n", pokemonResp.Height, pokemonResp.Weight)
	fmt.Println("\tStats:")
	for _, stat := range pokemonResp.Stats {
		fmt.Printf("\t\t󰝪 %v:\t%v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("\tTypes:")
	for _, t := range pokemonResp.Types {
		fmt.Printf("\t\t %v\n", t.Type.Name)
	}
	fmt.Println("\tAbilities")
	for _, a := range pokemonResp.Abilities {
		fmt.Printf("\t\t󰓥 %v\n", a.Ability.Name)
	}

	return nil
}
