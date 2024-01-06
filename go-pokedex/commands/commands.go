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
	}
}
