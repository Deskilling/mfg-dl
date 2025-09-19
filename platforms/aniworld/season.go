package aniworld

import (
	parser "mfg-dl/parser/aniworld"

	"fmt"
	"mfg-dl/request"

	"github.com/charmbracelet/log"
)

func Seasons(anime string) ([]parser.Season, error) {
	seasons, err := request.Get(request.AniworldEndpoints["episodes"] + anime)
	if err != nil {
		err = fmt.Errorf("failed to GET Seasons for %s with error %w", anime, err)
		log.Error(err)
		return nil, err
	}

	parsedSeasons, err := parser.Seasons(seasons)
	if err != nil {
		err = fmt.Errorf("failed parsing seasons for %s: %w", anime, err)
		log.Error(err)
		return nil, err
	}
	if len(parsedSeasons) == 0 {
		err = fmt.Errorf("%s not found", anime)
		log.Error(err)
		return nil, err
	}

	return parsedSeasons, nil
}
