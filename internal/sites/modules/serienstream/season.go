package serienstream

import (
	"fmt"
	"strings"

	"mfg-dl/internal/sites/model"

	"github.com/gocolly/colly/v2"
)

func (service *Serienstream) Seasons(result model.SearchResult) (seasons []model.Season, err error) {
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	c.OnHTML("ul.nav.list-items-nav li.nav-item a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if href == "" {
			return
		}

		if strings.Contains(href, "/episode-") {
			return
		}

		label := e.Text
		if label == "" {
			return
		}

		isFilme := strings.Contains(href, "/filme") || strings.Contains(href, "/staffel-0")
		isSeason := strings.Contains(href, "/staffel-") && !strings.Contains(href, "/staffel-0")

		if !isFilme && !isSeason {
			return
		}

		num := "00"
		if isSeason {
			num = fmt.Sprintf("%02s", strings.TrimSpace(label))
			label = "Season" + label
		}

		seasons = append(seasons, model.Season{
			Service:     Name,
			Name:        result.Name,
			Href:        href,
			SeasonNum:   num,
			SeasonLabel: strings.TrimSpace(label),
		})
	})

	err = c.Visit(BaseURL + result.Href)
	if err != nil {
		fmt.Printf("Failed to visit initial URL: %v\n", err)
	}
	c.Wait()

	return seasons, nil
}
