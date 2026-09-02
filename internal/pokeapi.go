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

func Pokeapi(url string, poke_pointer *Pokestruct, poke_cache *Cache) error {

	body, ok := poke_cache.Get(url)

	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d", res.StatusCode)
		}

		poke_cache.Add(url, body)
	}

	j_err := json.Unmarshal(body, poke_pointer)
	if j_err != nil {
		return j_err
	}
	return nil
}
