package aniworldParser

import (
	"encoding/json"
	"fmt"
	"html"

	"github.com/charmbracelet/log"
)

type SearchResult struct {
	Name           string `json:"name"`
	Link           string `json:"link"`
	Description    string `json:"description"`
	Cover          string `json:"cover"`
	ProductionYear string `json:"productionYear"`
}

func Search(data string) ([]SearchResult, error) {
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
