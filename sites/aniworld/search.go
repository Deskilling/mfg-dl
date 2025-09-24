package aniworld

import (
	"encoding/json"
	"html"

	"fmt"
	"mfg-dl/request"
	"mfg-dl/util"

	"github.com/charmbracelet/log"
)

type SearchResult struct {
	Name           string `json:"name"`
	Link           string `json:"link"`
	Description    string `json:"description"`
	Cover          string `json:"cover"`
	ProductionYear string `json:"productionYear"`
}

func GetSearch(term string) ([]SearchResult, error) {
	encodedTerm := util.EncodeURIComponent(term)

	searchResults, err := request.Get(request.AniworldEndpoints["search"] + encodedTerm)
	if err != nil {
		err = fmt.Errorf("failed to GET Search for %s: %w", term, err)
		log.Error(err)
		return nil, err
	}

	parsedResults, err := parseSearch(searchResults)
	if err != nil {
		err = fmt.Errorf("failed parsing search results for %s: %w", term, err)
		log.Error(err)
		return nil, err
	}
	if len(parsedResults) == 0 {
		err = fmt.Errorf("%s not found", term)
		log.Error(err)
		return nil, err
	}

	return parsedResults, nil
}

func parseSearch(data string) ([]SearchResult, error) {
	var search []SearchResult

	err := json.Unmarshal([]byte(data), &search)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal search results: %w", err)
		log.Error(err)
		return nil, err
	}

	for i := range search {
		search[i].Name = html.UnescapeString(search[i].Name)
		search[i].Description = html.UnescapeString(search[i].Description)
	}

	return search, nil
}
