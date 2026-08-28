package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type pokestruct struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []city `json:"results"`
}

type city struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func pokeapi(url string, poke_pointer *pokestruct) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d", res.StatusCode)
	}

	j_err := json.Unmarshal(body, poke_pointer)
	if j_err != nil {
		return j_err
	}
	return nil
}
