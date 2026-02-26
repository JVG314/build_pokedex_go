package pokeapi

import (
	"encoding/json"
	"net/http"
	"time"
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

type Client struct {
	httpClient http.Client
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) GetLocationAreas(url string) (RespLocationAreas, error) {
	if url == "" {
		url = baseURL + "/location-area"
	}
	res, err := c.httpClient.Get(url)
	if err != nil {
		return RespLocationAreas{}, err
	}
	defer res.Body.Close()

	var respLA RespLocationAreas
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&respLA)
	if err != nil {
		return RespLocationAreas{}, err
	}
	return respLA, nil
}
