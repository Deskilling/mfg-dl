package aniworld

import (
	"encoding/json"
	"fmt"
	"html"
	"time"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

	"github.com/charmbracelet/log"
)

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

	if len(parsedResults) == 0 {
		err = fmt.Errorf("%s not found", term)
		log.Error(err)
		return nil, err
	}

	return parsedResults, nil
}

func ParseSearch(data string) (search []model.SearchResult, err error) {
	err = json.Unmarshal([]byte(data), &search)
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
