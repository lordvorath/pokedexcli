package poke_api

import (
	"encoding/json"
	"net/http"
)

type NamedAPIResource struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type NamedAPIResourceList struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []NamedAPIResource `json:"results"`
}

func GetLocationAreas(url string) (NamedAPIResourceList, error) {
	var mapList NamedAPIResourceList
	client := http.DefaultClient
	res, err := client.Get(url)
	if err != nil {
		return mapList, err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&mapList)
	if err != nil {
		return mapList, err
	}
	return mapList, nil
}
