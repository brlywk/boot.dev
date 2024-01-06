package commands

import (
	"brlywk/bootdev/pokedex/pokeapi"
	"os"
)

// Exit the Gokedex
func CommandExit(_ *pokeapi.ApiConfig, _ string) error {
	os.Exit(0)
	return nil
}
