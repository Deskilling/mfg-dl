package aniworld

import (
	"fmt"
	"strings"

	"mfg-dl/internal/sites/model"

	"github.com/gocolly/colly/v2"
)

func (service *Aniworld) Streams(episode model.Episode) (streams []model.Stream, err error) {
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	c.OnHTML("ul.row li", func(e *colly.HTMLElement) {
		hosterName := e.ChildText("h4")
		if hosterName == "" {
			hosterName = e.ChildAttr("i.icon", "title")
			hosterName = strings.TrimPrefix(hosterName, "Hoster ")
		}

		linkTarget := e.Attr("data-link-target")
		if linkTarget == "" {
			linkTarget = e.ChildAttr("a.watchEpisode", "href")
		}

		langKey := e.Attr("data-lang-key")

		streams = append(streams, model.Stream{
			Service:                 Name,
			Name:                    episode.Name,
			Href:                    linkTarget,
			SeasonNum:               episode.SeasonNum,
			EpisodeTitle:            episode.EpisodeTitle,
			EpisodeAlternativeTitle: episode.EpisodeAlternativeTitle,
			EpisodeNum:              episode.EpisodeNum,
			Hoster:                  hosterName,
			Language:                AniLanguages[langKey],
		})
	})

	err = c.Visit(BaseURL + episode.Href)
	if err != nil {
		fmt.Printf("Failed to visit initial URL: %v\n", err)
	}
	c.Wait()

	return streams, nil
}
