package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lordvorath/pokedexcli/internal/poke_api"
	"github.com/lordvorath/pokedexcli/internal/pokecache"
)

func main() {
	config := &poke_api.Config{
		Next:     "",
		Previous: "",
		Cache:    pokecache.NewCache(5 * time.Minute),
		Pokedex:  make(map[string]poke_api.Pokemon),
	}
	reader := bufio.NewReader(os.Stdin)
	input := bufio.NewScanner(reader)

	for {
		fmt.Print("Pokedex > ")
		var text string
		if !input.Scan() {
			os.Exit(0)
			return
		}
		text = input.Text()
		words := cleanInput(text)
		commandWord, args := words[0], words[1:]
		command, exists := getCommands()[commandWord]
		if exists {
			err := command.callback(config, args)
			if err != nil {
				fmt.Printf("Error while executing callback: %v\n", err)
				continue
			}
		} else {
			fmt.Print("Unknown command\n")
			continue
		}
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}
