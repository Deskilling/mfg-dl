package aniworld

import (
	"fmt"
	"log/slog"
	"mfg-dl/internal/sites/model"
	"strings"

	"github.com/gocolly/colly/v2"
)

func (service *Aniworld) Seasons(result model.SearchResult) (seasons []model.Season, err error) {
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	c.OnHTML(".hosterSiteDirectNav a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if href == "" {
			return
		}

		label := e.Attr("title")
		if label == "" {
			return
		}

		num := "0"
		isFilme := strings.Contains(href, "/filme")
		isSeason := strings.Contains(href, "/staffel-")

		if !isFilme && !isSeason {
			return
		}

		if isSeason {
			num = strings.TrimPrefix(label, "Staffel ")
		}

		formattedNumber := fmt.Sprintf("%02s", num)

		if strings.Contains(href, "/episode-") {
			return
		}

		seasons = append(seasons, model.Season{
			Service:     Name,
			Name:        result.Name,
			Href:        href,
			SeasonNum:   formattedNumber,
			SeasonLabel: label,
		})
	})

	url := BaseURL + result.Href
	err = c.Visit(url)
	if err != nil {
		slog.Error("failed visiting site", "url", url, "err", err)
	}
	c.Wait()

	return seasons, nil
}
