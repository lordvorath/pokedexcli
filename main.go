package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lordvorath/pokedexcli/internal/pokecache"
)

func main() {
	config := &Config{}
	reader := bufio.NewReader(os.Stdin)
	input := bufio.NewScanner(reader)
	cache := pokecache.NewCache(15 * time.Second)

	for {
		fmt.Print("Pokedex > ")
		var text string
		if !input.Scan() {
			os.Exit(0)
			return
		}
		text = input.Text()
		words := cleanInput(text)
		command, exists := getCommands()[words[0]]
		if exists {
			err := command.callback(config, cache)
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
