package pokeapi

import (
	"brlywk/bootdev/pokedex/cache"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Get list of areas
func FetchLocationData(url string, cache *cache.Cache) (LocationResponse, error) {
	return FetchData[LocationResponse](url, cache)
}

// Get exploration data for a single area
func FetchExplorationData(url string, area string, cache *cache.Cache) (ExplorationResponse, error) {
	locationUrl := fmt.Sprintf("%v/%v", url, area)

	return FetchData[ExplorationResponse](locationUrl, cache)
}

// Fetches Data and creates response struct
//
// Checks if url has already been fetched and returns fetched data if available
func FetchData[T any](url string, cache *cache.Cache) (T, error) {
	var resp T
	var err error

	body, found := cache.Get(url)

	if !found {
		res, err := http.Get(url)
		if err != nil {
			return resp, err
		}

		if res.StatusCode > 299 {
			return resp, fmt.Errorf("Response failed with code %v", res.StatusCode)
		}

		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return resp, err
		}

		cache.Add(url, body)
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
