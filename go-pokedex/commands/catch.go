package commands

import (
	"brlywk/bootdev/pokedex/pokeapi"
	"brlywk/bootdev/pokedex/utils"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Catch a pokemon
func CommandCatch(config *pokeapi.ApiConfig, pokemon string) error {
	if pokemon == "" {
		return errors.New("Please provide an area to explore")
	}

	resp, err := pokeapi.FetchPokemonData(config.PokemonUrl, pokemon, config.Cache)
	if err != nil {
		return err
	}

	fmt.Printf("\tThrowing a Pokeball at %v%v%v\n", utils.PromptColorOrange, resp.Name, utils.PromptColorReset)

	fmt.Print("\t")
	for i := 0; i < 5; i++ {
		fmt.Print(".")
		time.Sleep(200 * time.Millisecond)
	}

	exp := resp.BaseExperience
	rnd := rand.Intn(exp)

	if rnd > (exp / 2) {
		fmt.Printf("\t%vYou caught %v%v!\n", utils.PromptColorGreen, resp.Name, utils.PromptColorReset)
		caught := *config.CaughtPokemon
		caught[pokemon] = resp
	} else {
		fmt.Printf("\t%v%v escaped...%v\n", utils.PromptColorRed, resp.Name, utils.PromptColorReset)
	}

	return nil
}
