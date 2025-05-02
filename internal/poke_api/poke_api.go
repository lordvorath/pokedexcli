package poke_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lordvorath/pokedexcli/internal/pokecache"
)

type Config struct {
	client   http.Client
	Next     string
	Previous string
	Cache    *pokecache.Cache
	Pokedex  map[string]Pokemon
}

func makeAPIRequest(url string, config *Config) ([]byte, error) {
	if val, ok := config.Cache.Get(url); ok {
		// fmt.Println("Returning cached result")
		return val, nil
	}
	client := http.DefaultClient
	res, err := client.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	// fmt.Println("Got response from API")
	byteSlice, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	config.Cache.Add(url, byteSlice)
	return byteSlice, nil
}

func GetLocationAreas(url string, config *Config) (NamedAPIResourceList, error) {
	var mapList NamedAPIResourceList
	byteSlice, err := makeAPIRequest(url, config)
	if err != nil {
		return mapList, err
	}

	err = json.Unmarshal(byteSlice, &mapList)
	if err != nil {
		return mapList, fmt.Errorf("failed to unmarshal location list: %w", err)
	}
	return mapList, nil
}

func GetExploredArea(url string, config *Config) (LocationArea, error) {
	var exploredLocation LocationArea
	byteSlice, err := makeAPIRequest(url, config)
	if err != nil {
		return exploredLocation, err
	}
	err = json.Unmarshal(byteSlice, &exploredLocation)
	if err != nil {
		return exploredLocation, fmt.Errorf("failed to unmarshal location: %w", err)
	}
	return exploredLocation, nil
}

func GetPokemon(url string, config *Config) (Pokemon, error) {
	var pokemon Pokemon
	byteSlice, err := makeAPIRequest(url, config)
	if err != nil {
		return pokemon, err
	}
	err = json.Unmarshal(byteSlice, &pokemon)
	if err != nil {
		return pokemon, fmt.Errorf("failed to unmarshal pokemon: %w", err)
	}
	return pokemon, nil
}
