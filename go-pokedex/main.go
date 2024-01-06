package main

import (
	"brlywk/bootdev/pokedex/cache"
	"brlywk/bootdev/pokedex/commands"
	"brlywk/bootdev/pokedex/pokeapi"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	apiUrl      = "https://pokeapi.co/api/v2/location-area"
	promptColor = "\033[38;5;124m"
	promptReset = "\033[0m"
)

// ----- Main ------------------------------------

func main() {
	cache := cache.NewCache(20 * time.Second)
	go cache.CleanupLoop()

	config := pokeapi.ApiConfig{
		BaseUrl:          "https://pokeapi.co/api/v2",
		LocationPath:     "location-area",
		NextLocation:     apiUrl,
		PreviousLocation: "",
		Cache:            &cache,
	}
	commandList := commands.GetCommands()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%vpokedex > %v", promptColor, promptReset)
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
