package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

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
			description: "\t\t\tDisplays a help message",
			callback:    commandHelp,
		},
		"catch": {
			name:        "catch <pokemon>",
			description: "\tAttempt to catch a Pokemon and add it to your Pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon>",
			description: "\tShow details of a Pokemon in your Pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "\t\tList captured Pokemon that can be inspected",
			callback:    commandPokedex,
		},
		"exit": {
			name:        "exit",
			description: "\t\t\tExit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore <area_name>",
			description: "\tLists pokemon found in the chosen area",
			callback:    commandExplore,
		},
		"map": {
			name:        "map",
			description: "\t\t\tDisplays the next 20 locations in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "\t\t\tDisplays the previous 20 locations in the Pokemon world",
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
	exploredLocation, err := poke_api.GetExploredArea(baseUrl+args[0], config)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", args[0])
	fmt.Println("Found Pokemon: ")
	var pokemonList []string
	for _, val := range exploredLocation.PokemonEncounters {
		pokemonList = append(pokemonList, val.Pokemon.Name)
	}
	for _, pokemon := range pokemonList {
		fmt.Println(pokemon)
	}
	return nil
}

func commandCatch(config *poke_api.Config, args []string) error {
	if len(args) <= 0 || len(args) > 1 {
		return fmt.Errorf("missing argument <pokemon> or too many arguments")
	}
	if _, ok := config.Pokedex[args[0]]; ok {
		return fmt.Errorf("you have already caught that pokemon")
	}
	baseUrl := "https://pokeapi.co/api/v2/pokemon/"
	pokemon, err := poke_api.GetPokemon(baseUrl+args[0], config)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)
	if caught := catchPokemon(pokemon.BaseExperience); caught {
		fmt.Printf("%v was caught!\n", pokemon.Name)
		config.Pokedex[pokemon.Name] = pokemon
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%v escaped!\n", pokemon.Name)
	}
	return nil
}

func catchPokemon(BaseExperience int) bool {
	rng := rand.New(rand.NewSource(int64(time.Now().Unix())))
	n := rng.Intn(BaseExperience)
	catchThreshold := 5000 // TODO put back to 50 or so
	// fmt.Printf("Rolled %d out of %d\n", n, catchThreshold)
	return n <= catchThreshold
}

func commandInspect(config *poke_api.Config, args []string) error {
	if pokemon, ok := config.Pokedex[args[0]]; !ok {
		return fmt.Errorf("you have not caught that pokemon")
	} else {
		fmt.Printf("Name: %v\n", pokemon.Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)
		fmt.Printf("Stats:\n")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%v: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Printf("Types:\n")
		for _, t := range pokemon.Types {
			fmt.Printf("  -%v\n", t.Type.Name)
		}
		return nil
	}
}

func commandPokedex(config *poke_api.Config, args []string) error {
	if len(config.Pokedex) < 1 {
		return fmt.Errorf("you haven't caught any Pokemon")
	}
	fmt.Println("Your Pokedex:")
	for _, pokemon := range config.Pokedex {
		fmt.Printf("  -%s\n", pokemon.Name)
	}
	return nil
}
