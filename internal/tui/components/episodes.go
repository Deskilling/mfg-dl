package components

import (
	"errors"
	"fmt"
	"strings"

	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

	"charm.land/log/v2"
)

func Episodes(site model.Site, season model.Season) (streams []model.Stream, err error) {
	episodes, err := site.Episodes(season)
	if err != nil {
		log.Error("Failed getting episodes", "service", site.Name(), "href", season.Href)
		return []model.Stream{}, err
	}

	if len(episodes) == 0 {
		return []model.Stream{}, errors.New("no episodes found")
	}

	var u []int
	if len(episodes) == 1 {
		log.Info("Selected the only episode", "episodeNum", episodes[0].EpisodeNum, "seasonNum", season.SeasonNum)
		u = append(u, 0)
	} else {
		for {
			for i := range episodes {
				if strings.TrimSpace(episodes[i].EpisodeTitle) == "" {
					episodes[i].EpisodeTitle = episodes[i].EpisodeAlternativeTitle
				}

				fmt.Printf("[%v] %s\n", i+1, episodes[i].EpisodeTitle)
			}

			v, err := ReadString(Reader, "Select (use 1,2,3,4 or all for selection): ")
			if err != nil {
				return []model.Stream{}, fmt.Errorf("failed userinput: %s", err)
			}

			if v == "all" {
				u = make([]int, len(episodes))
				for i := range episodes {
					u[i] = i
				}
				break
			}

			valid := []int{}
			for _, w := range util.ParseMultipleInts(v) {
				w--
				if w >= 0 && w < len(episodes) {
					valid = append(valid, w)
				}
			}

			if len(valid) > 0 {
				u = valid
				break
			}
		}
	}

	lang, err := Language(site, episodes[u[0]])
	if err != nil {
		return []model.Stream{}, fmt.Errorf("failed getting language: %w", err)
	}

	for _, v := range u {
		stream, err := site.Streams(episodes[v])
		if err != nil {
			continue
		}

		for _, w := range stream {
			if w.Language == lang {
				streams = append(streams, w)
				break
			}
		}
	}

	return streams, nil
}
