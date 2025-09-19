package aniworld

import (
	parser "mfg-dl/parser/aniworld"

	"fmt"
	"mfg-dl/request"

	"github.com/charmbracelet/log"
)

func Episodes(anime, season string) ([]parser.Episode, error) {
	episodes, err := request.Get(request.AniworldEndpoints["episodes"] + anime + "/staffel-" + season)
	if err != nil {
		err = fmt.Errorf("failed to GET Episodes for %s: %w", anime, err)
		log.Error(err)
		return nil, err
	}

	parsedEpisodes, err := parser.Episodes(episodes)
	if err != nil {
		err = fmt.Errorf("failed parsing episodes for %s %s: %w", anime, season, err)
		log.Error(err)
		return nil, err
	}
	if len(parsedEpisodes) == 0 {
		err = fmt.Errorf("%s %s not found", anime, season)
		log.Error(err)
		return nil, err
	}

	return parsedEpisodes, nil
}
