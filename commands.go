package main

import (
	"fmt"
	"os"

	"github.com/lordvorath/pokedexcli/internal/poke_api"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *poke_api.Config, args []string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore <area_name>",
			description: "Lists pokemon found in the chosen area",
			callback:    commandExplore,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations in the Pokemon world",
			callback:    commandMapb,
		},
	}
}

func commandExit(config *poke_api.Config, args []string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("dalfhjlksdjflakhd")
}

func commandHelp(config *poke_api.Config, args []string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, info := range getCommands() {
		fmt.Printf("%v: %v\n", info.name, info.description)
	}
	return nil
}

func commandMap(config *poke_api.Config, args []string) error {
	var url string
	if config.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = config.Next
	}
	// fmt.Println("Calling GetLocationAreas")
	mapList, err := poke_api.GetLocationAreas(url, config)
	if err != nil {
		return err
	}
	// fmt.Println("Received mapList")
	if mapList.Next != nil {
		config.Next = *mapList.Next
	}
	if mapList.Previous != nil {
		config.Previous = *mapList.Previous
	}
	for _, location := range mapList.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapb(config *poke_api.Config, args []string) error {
	var url string
	if config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		url = config.Previous
	}
	mapList, err := poke_api.GetLocationAreas(url, config)
	if err != nil {
		return err
	}
	if mapList.Next != nil {
		config.Next = *mapList.Next
	}
	if mapList.Previous != nil {
		config.Previous = *mapList.Previous
	}
	for _, location := range mapList.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandExplore(config *poke_api.Config, args []string) error {
	if len(args) <= 0 || len(args) > 1 {
		return fmt.Errorf("missing argument <area_name> or too many arguments")
	}
	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	pokemonList, err := poke_api.GetExploredArea(baseUrl+args[0], config)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", args[0])
	fmt.Println("Found Pokemon: ")
	for _, pokemon := range pokemonList {
		fmt.Println(pokemon)
	}
	return nil
}
