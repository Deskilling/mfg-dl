package aniworld

import (
	"encoding/json"
	"fmt"
	"html"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

	"charm.land/log/v2"
)

type SearchResult struct {
	Name           string `json:"name"`
	Href           string `json:"link"`
	Description    string `json:"description"`
	Cover          string `json:"cover"`
	ProductionYear string `json:"productionYear"`
}

func GetSearch(term string) ([]model.SearchResult, error) {
	encodedTerm := util.EncodeURIComponent(term)

	searchResults, err := request.Get(AniEndpoints["search"] + encodedTerm)
	if err != nil {
		log.Error("failed to get search", "service", Name, "term", term)
		return nil, fmt.Errorf("failed to GET Search for %s: %w", term, err)
	}

	parsedResults, err := ParseSearch(searchResults)
	if err != nil {
		log.Error("failed to parse search", "service", Name, "term", term)
		return nil, fmt.Errorf("failed parsing search results for %s: %w", term, err)
	}

	return parsedResults, nil
}

func ParseSearch(data []byte) (search []model.SearchResult, err error) {
	var searchResults []SearchResult

	err = json.Unmarshal(data, &searchResults)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal search results: %w", err)
		return nil, err
	}

	for _, v := range searchResults {
		result := model.SearchResult{
			Service:        Name,
			Name:           html.UnescapeString(v.Name),
			Href:           v.Href,
			Description:    html.UnescapeString(v.Description),
			Cover:          v.Cover,
			ProductionYear: v.ProductionYear,
		}

		search = append(search, result)
	}

	return search, nil
}
