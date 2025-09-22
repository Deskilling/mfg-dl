package aniworld

import (
	aniworldParser "mfg-dl/sites/aniworld/parser"

	"fmt"
	"mfg-dl/request"
	"mfg-dl/util"

	"github.com/charmbracelet/log"
)

func Search(term string) ([]aniworldParser.SearchResult, error) {
	encodedTerm := util.EncodeURIComponent(term)

	searchResults, err := request.Get(request.AniworldEndpoints["search"] + encodedTerm)
	if err != nil {
		err = fmt.Errorf("failed to GET Search for %s: %w", term, err)
		log.Error(err)
		return nil, err
	}

	parsedResults, err := aniworldParser.Search(searchResults)
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
