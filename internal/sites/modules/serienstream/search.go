package serienstream

import (
	"encoding/json"
	"fmt"
	"html"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"
)

type SearchResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type searchResponse struct {
	Shows []SearchResult `json:"shows"`
}

func (service *Serienstream) Search(term string) ([]model.SearchResult, error) {
	encodedTerm := util.EncodeURIComponent(term)

	searchResults, err := request.Get(service.client, SerienstreamEndpoints["search"]+encodedTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to GET Search for %s: %w", term, err)
	}

	parsedResults, err := ParseSearch(searchResults)
	if err != nil {
		return nil, fmt.Errorf("failed parsing search results: %w", err)
	}

	return parsedResults, nil
}

func ParseSearch(data []byte) (search []model.SearchResult, err error) {
	var response searchResponse

	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	for _, v := range response.Shows {
		search = append(search, model.SearchResult{
			Service: Name,
			Name:    html.UnescapeString(v.Name),
			Href:    v.URL,
		})
	}

	return search, nil
}
