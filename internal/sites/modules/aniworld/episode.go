package aniworld

import (
	"fmt"
	"log/slog"
	"mfg-dl/internal/sites/model"

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

	url := BaseURL + season.Href
	err = c.Visit(url)
	if err != nil {
		slog.Error("failed visiting site", "url", url, "err", err)
	}
	c.Wait()

	return episodes, err
}
