package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Print("Hello, World!")
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}
