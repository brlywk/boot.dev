package commands

import "brlywk/bootdev/pokedex/pokeapi"

// ----- Structs ---------------------------------

// Commands executable by the user
type CliCommand struct {
	Name        string
	Description string
	Callback    func(config *pokeapi.ApiConfig, param string) error
}

// Returns a map of all available commands
func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    CommandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Gokedex",
			Callback:    CommandExit,
		},
		"map": {
			Name:        "map",
			Description: "Displays the next 20 pokemon locations",
			Callback:    CommandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays the previous 20 pokemon locations",
			Callback:    CommandMapB,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore an area.\n\tUse 'explore <area name>' to select which area.\n\n\tExample:\texplore canalave-city",
			Callback:    CommandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Attempt to catch a Pokemon.\n\tUse 'catch <name>' to select Pokemon.\n\n\tExample:\tcatch pikachu",
			Callback:    CommandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect a Pokemon.\n\tUse 'inspect <name>' to select Pokemon.\n\n\tExample:\tinspect pikachu",
			Callback:    CommandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "View your current Pokedex (all Pokemon you have caught so far.)",
			Callback:    CommandPokedex,
		},
	}
}
