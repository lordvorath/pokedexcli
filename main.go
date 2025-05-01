package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	input := bufio.NewScanner(reader)

	for {
		fmt.Print("Pokedex > ")
		var text string
		if input.Scan() {
			text = input.Text()
			words := cleanInput(text)
			fmt.Printf("Your command was: %s\n", words[0])
		} else {
			return
		}
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}
