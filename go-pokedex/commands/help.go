package commands

import (
	"brlywk/bootdev/pokedex/pokeapi"
	"fmt"
)

// Display list of commands and descriptions
func CommandHelp(_ *pokeapi.ApiConfig, _ string) error {
	fmt.Println("Welcome to the Gokedex!")
	fmt.Printf("\nUsage:\n\n")

	for _, info := range GetCommands() {
		fmt.Printf("%s\t%s\n", info.Name, info.Description)
	}

	fmt.Println("")

	return nil
}
