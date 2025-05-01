package main

import (
	"fmt"
	"os"

	"internal/poke_api"
)

type Config struct {
	Next     string
	Previous string
}
type cliCommand struct {
	name        string
	description string
	callback    func(config *Config) error
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

func commandExit(config *Config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("dalfhjlksdjflakhd")
}

func commandHelp(config *Config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, info := range getCommands() {
		fmt.Printf("%v: %v\n", info.name, info.description)
	}
	return nil
}

func commandMap(config *Config) error {
	var url string
	if config.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = config.Next
	}
	mapList, err := poke_api.GetLocationAreas(url)
	if err != nil {
		return err
	}
	config.Next = mapList.Next
	config.Previous = mapList.Previous
	for _, location := range mapList.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapb(config *Config) error {
	var url string
	if config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		url = config.Previous
	}
	mapList, err := poke_api.GetLocationAreas(url)
	if err != nil {
		return err
	}
	config.Next = mapList.Next
	config.Previous = mapList.Previous
	for _, location := range mapList.Results {
		fmt.Println(location.Name)
	}
	return nil
}
