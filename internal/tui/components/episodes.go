package components

import (
	"fmt"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"
	"strings"
)

func Episodes(site *model.Site, season model.Season) []model.Stream {
	episodes, err := site.Episodes(season)
	if err != nil {
		return nil
	}

	for i := range episodes {
		if strings.TrimSpace(episodes[i].EpisodeTitle) == "" {
			episodes[i].EpisodeTitle = episodes[i].EpisodeAlternativeTitle
		}
	}

	var u []int

	if len(episodes) == 1 {
		u = append(u, 0)
	} else {
		for {
			for i := range episodes {
				fmt.Printf("[%v] %s\n", i+1, episodes[i].EpisodeTitle)
			}

			v := ReadString(Reader, "Select (use 1,2,3,4 or all for selection): ")

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
		return nil
	}

	var unc []model.Stream
	for _, v := range u {
		stream, err := site.Streams(episodes[v])
		if err != nil {
			continue
		}

		for _, w := range stream {
			if w.Language == lang {
				unc = append(unc, w)
				break
			}
		}
	}

	return unc
}
