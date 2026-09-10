package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	HttpClient *http.Client
	Cache      *Cache
}

type Pokestruct struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []city `json:"results"`
}

type city struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Cityinfo struct {
	Id         int          `json:"id"`
	Name       string       `json:"name"`
	Game_index int          `json:"game_index"`
	Pkmn_enc   []Encounters `json:"pokemon_encounters"`
}

type Encounters struct {
	Pokemon struct {
		Name     string `json:"name"`
		Pkmn_url string `json:"url"`
	}
}

type Pokemon struct {
	Name     string      `json:"name"`
	Height   int         `json:"height"`
	Weight   int         `json:"weight"`
	Base_exp int         `json:"base_experience"`
	Stats    []Stats     `json:"stats"`
	Types    []Poketypes `json:"types"`
}

type Stats struct {
	Basestat int  `json:"base_stat"`
	Stat     Stat `json:"stat"`
}

type Stat struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Poketypes struct {
	Slot     int      `json:"slot"`
	Poketype Poketype `json:"type"`
}

type Poketype struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (c Client) ListLocations(url string, poke_pointer *Pokestruct, poke_cache *Cache) error {
	body, ok := poke_cache.Get(url)

	var err error
	if !ok {
		body, err = getAPIData(url)
		poke_cache.Add(url, body)
	}
	if err != nil {
		return err
	}

	j_err := json.Unmarshal(body, poke_pointer)
	if j_err != nil {
		return j_err
	}
	return nil
}

func (c Client) GetLocationPkmn(location string, city_pointer *Cityinfo, poke_cache *Cache) error {
	url := "https://pokeapi.co/api/v2/location-area/" + location

	body, ok := poke_cache.Get(url)

	var err error
	if !ok {
		body, err = getAPIData(url)
		poke_cache.Add(url, body)
	}
	if err != nil {
		return err
	}

	j_err := json.Unmarshal(body, city_pointer)
	if j_err != nil {
		return j_err
	}
	return nil
}

func (c Client) GetPokemoninfo(poke_string string, pokemon_pointer *Pokemon, poke_cache *Cache) error {
	url := "https://pokeapi.co/api/v2/pokemon/" + poke_string

	body, ok := poke_cache.Get(url)

	var err error
	if !ok {
		body, err = getAPIData(url)
		poke_cache.Add(url, body)
	}
	if err != nil {
		return err
	}

	j_err := json.Unmarshal(body, pokemon_pointer)
	if j_err != nil {
		return j_err
	}

	return nil
}

func getAPIData(url string) ([]byte, error) {
	var blank_body []byte
	res, err := http.Get(url)
	if err != nil {
		return blank_body, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return blank_body, err
	}

	if res.StatusCode > 299 {
		return body, fmt.Errorf("Response failed with status code: %d", res.StatusCode)
	}

	return body, err
}
