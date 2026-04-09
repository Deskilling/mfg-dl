package components

import (
	"errors"
	"fmt"
	"strings"

	"mfg-dl/internal/sites/model"

	"charm.land/huh/v2"
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
		var values []huh.Option[int]
		for i, v := range episodes {
			title := strings.TrimSpace(v.EpisodeTitle)
			if title == "" {
				title = v.EpisodeAlternativeTitle
			}

			title = fmt.Sprintf("%v %s", v.EpisodeNum, title)

			values = append(values, huh.NewOption(title, i))
		}

		var selected []int
		err := huh.NewMultiSelect[int]().
			Title("Select episodes").
			Options(values...).
			Value(&selected).
			WithTheme(Theme).
			Run()
		if err != nil {
			return []model.Stream{}, fmt.Errorf("failed userinput: %w", err)
		}
		u = selected
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
