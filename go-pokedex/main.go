package main

import (
	"brlywk/bootdev/pokedex/cache"
	"brlywk/bootdev/pokedex/commands"
	"brlywk/bootdev/pokedex/pokeapi"
	"brlywk/bootdev/pokedex/utils"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// ----- Main ------------------------------------

func main() {
	// Initialise cache and start cleanup loop
	cache := cache.NewCache(20 * time.Second)
	go cache.CleanupLoop()

	config := pokeapi.ApiConfig{
		LocationAreaUrl:  "https://pokeapi.co/api/v2/location-area",
		PokemonUrl:       "https://pokeapi.co/api/v2/pokemon",
		NextLocation:     "https://pokeapi.co/api/v2/location-area",
		PreviousLocation: "",
		Cache:            &cache,
		CaughtPokemon:    &map[string]pokeapi.PokemonResponse{},
	}

	commandList := commands.GetCommands()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%vpokedex > %v", utils.PromptColorOrange, utils.PromptColorReset)
		scanner.Scan()

		input := scanner.Text()
		cmd, param := parseInput(input)
		command, ok := commandList[cmd]

		if ok {
			err := command.Callback(&config, param)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// Parse user input and split into command and parameter
func parseInput(input string) (string, string) {
	fields := strings.Fields(strings.ToLower(input))

	l := len(fields)

	switch {
	case l == 1:
		return fields[0], ""
	case l > 1:
		return fields[0], fields[1]
	default:
		return "", ""
	}
}
