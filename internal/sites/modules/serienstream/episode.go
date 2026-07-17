package serienstream

import (
	"fmt"
	"strings"

	"mfg-dl/internal/sites/model"

	"charm.land/log/v2"
	"github.com/gocolly/colly/v2"
)

func (service *Serienstream) Episodes(season model.Season) (episodes []model.Episode, err error) {
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	c.OnHTML("table.episode-table tbody tr.episode-row", func(e *colly.HTMLElement) {
		episodeNumberText := e.ChildText("th.episode-number-cell")
		if episodeNumberText == "" {
			return
		}

		formattedNumber := fmt.Sprintf("%02s", episodeNumberText)

		mainTitle := e.ChildText("td.episode-title-cell strong.episode-title-ger")
		secondaryTitle := e.ChildText("td.episode-title-cell span.episode-title-eng")

		onclickAttr := e.Attr("onclick")
		href := strings.TrimSuffix(strings.TrimPrefix(onclickAttr, "window.location='"), "'")
		if href == "" {
			return
		}

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
