package aniworld

import (
	parser "mfg-dl/parser/aniworld"

	"fmt"
	"mfg-dl/request"

	"github.com/charmbracelet/log"
)

func Streams(anime, season, episode string) ([]parser.Stream, error) {
	pageURL := request.AniworldEndpoints["episodes"] + anime + "/staffel-" + season + "/episode-" + episode
	log.Debug(pageURL)
	streams, err := request.Get(pageURL)
	if err != nil {
		err = fmt.Errorf("failed to GET Stream for %s %s %s: %w", anime, season, episode, err)
		log.Error(err)
		return nil, err
	}

	parsedStreams, err := parser.Streams(streams)
	if err != nil {
		err = fmt.Errorf("failed parsing Streams for %s %s %s: %w", anime, season, episode, err)
		log.Error(err)
		return nil, err
	}
	log.Debugf("parsed %d streams for %s %s %s", len(parsedStreams), anime, season, episode)
	if len(parsedStreams) == 0 {
		err = fmt.Errorf("%s not found", anime)
		log.Error(err)
		return nil, err
	}

	return parsedStreams, nil
}
