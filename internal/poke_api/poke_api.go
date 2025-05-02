package poke_api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/lordvorath/pokedexcli/internal/pokecache"
)

type Config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
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
		return mapList, err
	}
	return mapList, nil
}

func GetExploredArea(url string, config *Config) ([]string, error) {
	var pokemonList []string
	var exploredLocation LocationArea
	byteSlice, err := makeAPIRequest(url, config)
	if err != nil {
		return pokemonList, err
	}
	err = json.Unmarshal(byteSlice, &exploredLocation)
	if err != nil {
		return pokemonList, err
	}
	for _, val := range exploredLocation.PokemonEncounters {
		pokemonList = append(pokemonList, val.Pokemon.Name)
	}
	return pokemonList, nil
}
