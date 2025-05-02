package poke_api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/lordvorath/pokedexcli/internal/pokecache"
)

type NamedAPIResource struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type NamedAPIResourceList struct {
	Count    int                `json:"count"`
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
	Results  []NamedAPIResource `json:"results"`
}

func GetLocationAreas(url string, cache *pokecache.Cache) (NamedAPIResourceList, error) {
	// fmt.Println("Executing GetLocationAreas")
	var mapList NamedAPIResourceList
	if val, ok := cache.Get(url); ok {
		// fmt.Println("Returning cached result")
		err := json.Unmarshal(val, &mapList)
		if err != nil {
			return mapList, err
		}
		return mapList, nil
	}
	client := http.DefaultClient
	res, err := client.Get(url)
	if err != nil {
		return mapList, err
	}
	defer res.Body.Close()
	// fmt.Println("Got response from API")
	byteSlice, err := io.ReadAll(res.Body)
	if err != nil {
		return mapList, err
	}
	// fmt.Println("ReadAll of res.Body")
	cache.Add(url, byteSlice)
	err = json.Unmarshal(byteSlice, &mapList)
	if err != nil {
		return mapList, err
	}
	// fmt.Println("Unmarshalled JSON")
	return mapList, nil
}
