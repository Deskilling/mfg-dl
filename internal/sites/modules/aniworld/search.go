package aniworld

import (
	"encoding/json"
	"fmt"
	"html"
	"time"

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
		err = fmt.Errorf("failed to GET Search for %s: %w", term, err)
		log.Error(err)
		return nil, err
	}

	start := time.Now()
	parsedResults, err := ParseSearch(searchResults)
	if err != nil {
		err = fmt.Errorf("failed parsing search results for %s: %w", term, err)
		log.Error(err)
		return nil, err
	}
	log.Debugf("time took for search parsing: %v", time.Since(start))

	return parsedResults, nil
}

func ParseSearch(data string) (search []model.SearchResult, err error) {
	var searchResults []SearchResult

	err = json.Unmarshal([]byte(data), &searchResults)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal search results: %w", err)
		log.Error(err)
		return nil, err
	}

	for _, v := range searchResults {
		result := model.SearchResult{
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
