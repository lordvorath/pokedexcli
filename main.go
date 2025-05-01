package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var config *Config
	config = &Config{}
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
		command, exists := getCommands()[words[0]]
		if exists {
			err := command.callback(config)
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
