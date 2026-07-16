package aniworld

import (
	"fmt"
	"mfg-dl/internal/sites/model"

	"charm.land/log/v2"
	"github.com/gocolly/colly/v2"
)

func (service *Aniworld) Episodes(season model.Season) (episodes []model.Episode, err error) {
	c := colly.NewCollector(
		colly.MaxDepth(1),
	)

	c.OnHTML("table.seasonEpisodesList tbody tr", func(e *colly.HTMLElement) {
		numberText := e.Attr("data-episode-season-id")
		if numberText == "" {
			return
		}

		var numberStr string
		for _, char := range numberText {
			if char >= '0' && char <= '9' {
				numberStr += string(char)
			}
		}

		formattedNumber := fmt.Sprintf("%02s", numberStr)
		mainTitle := e.ChildText("td.seasonEpisodeTitle a strong")
		secondaryTitle := e.ChildText("td.seasonEpisodeTitle a span")

		href := e.ChildAttr("td.seasonEpisodeTitle a", "href")

		episodes = append(episodes, model.Episode{
			Service:                 Name,
			Name:                    season.Name,
			Href:                    href,
			SeasonNum:               season.SeasonNum,
			EpisodeTitle:            mainTitle,
			EpisodeAlternativeTitle: secondaryTitle,
			EpisodeNum:              formattedNumber,
		})
	})

	err = c.Visit(BaseURL + season.Href)
	if err != nil {
		log.Error(err)
	}
	c.Wait()

	return episodes, err
}
