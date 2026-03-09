package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/JVG314/build_pokedex_go/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type RespLocationAreas struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

type RespLocationAreaDetails struct {
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(5 * time.Second),
	}
}

func (c *Client) GetLocationAreas(url string) (RespLocationAreas, error) {
	if url == "" {
		url = baseURL + "/location-area/"
	}

	// Try cache first
	if cacheData, ok := c.cache.Get(url); ok {
		fmt.Printf("URL %s in cache\n", url)
		var respLA RespLocationAreas
		if err := json.Unmarshal(cacheData, &respLA); err != nil {
			return RespLocationAreas{}, err
		}
		return respLA, nil
	}

	// If not cached, do HTTP request
	fmt.Printf("URL %s not in cache, making HTTP request...\n", url)
	res, err := c.httpClient.Get(url)
	if err != nil {
		return RespLocationAreas{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return RespLocationAreas{}, fmt.Errorf("pokeapi: %s: %v", res.Status, strings.TrimSpace(string(body)))
	}

	// Read full body to cache it
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return RespLocationAreas{}, err
	}

	// Save to cache
	c.cache.Add(url, body)

	// Unmarshal and return
	var respLA RespLocationAreas
	if err := json.Unmarshal(body, &respLA); err != nil {
		return RespLocationAreas{}, err
	}

	return respLA, nil
}

func (c *Client) GetLocationDetails(name string) (RespLocationAreaDetails, error) {
	url := baseURL + "/location-area/" + name

	// Try cache first
	if cacheData, ok := c.cache.Get(url); ok {
		fmt.Printf("URL %s in cache\n", url)
		var respLAD RespLocationAreaDetails
		if err := json.Unmarshal(cacheData, &respLAD); err != nil {
			return RespLocationAreaDetails{}, err
		}
		return respLAD, nil
	}

	// If not cached, do HTTP request
	fmt.Printf("URL %s not in cache, making HTTP request...\n", url)
	res, err := c.httpClient.Get(url)
	if err != nil {
		return RespLocationAreaDetails{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return RespLocationAreaDetails{}, fmt.Errorf("pokeapi: %s: %v", res.Status, strings.TrimSpace(string(body)))
	}

	// Read full body to cache it
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return RespLocationAreaDetails{}, err
	}

	// Save to cache
	c.cache.Add(url, body)

	// Unmarshal and return
	var respLAD RespLocationAreaDetails
	if err := json.Unmarshal(body, &respLAD); err != nil {
		return RespLocationAreaDetails{}, err
	}

	return respLAD, nil
}

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + name
	res, err := c.httpClient.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return Pokemon{}, fmt.Errorf("pokeapi: %s: %v", res.Status, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	var respPokemon Pokemon
	if err := json.Unmarshal(body, &respPokemon); err != nil {
		return Pokemon{}, err
	}

	return respPokemon, nil
}
