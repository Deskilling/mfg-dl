package serienstream

import (
	"fmt"
	"strings"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"

	"github.com/PuerkitoBio/goquery"
)

type Season struct {
	Name        string
	Href        string
	SeasonNum   string
	SeasonLabel string
}

func (service *Serienstream) Seasons(result model.SearchResult) (seasons []model.Season, err error) {
	unparsedSeasons, err := request.Get(service.client, SerienstreamEndpoints["default"]+result.Href)
	if err != nil {
		return nil, err
	}

	parsedSeasons, err := ParseSeasons(string(unparsedSeasons))
	if err != nil {
		return nil, err
	}

	if len(parsedSeasons) == 0 {
		return nil, fmt.Errorf("%s not found", result.Href)
	}

	for i := range parsedSeasons {
		parsedSeasons[i].Name = result.Name
	}

	for _, v := range parsedSeasons {
		seasons = append(seasons, model.Season{
			Name:        v.Name,
			Href:        v.Href,
			SeasonNum:   v.SeasonNum,
			SeasonLabel: v.SeasonLabel,
		})
	}

	return seasons, nil
}

func ParseSeasons(html string) (seasons []Season, err error) {
	if html == "" {
		return nil, fmt.Errorf("empty html")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		return nil, fmt.Errorf("could not create goquery document: %w", err)
	}

	doc.Find("#season-nav a.alphabet-link").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		data, exists := s.Attr("data-season-pill")
		if !exists {
			return
		}

		label := strings.TrimSpace(s.Text())

		if !strings.Contains(href, "/staffel-") {
			return
		}

		var seasonNum string
		if data == "0" {
			seasonNum = "00"
			label = "Filme"
		} else {
			seasonNum = data
			if len(seasonNum) == 1 {
				seasonNum = "0" + seasonNum
			}
			label = "Staffel " + data
		}

		seasons = append(seasons, Season{
			Href:        href,
			SeasonLabel: label,
			SeasonNum:   seasonNum,
		})
	})

	return seasons, nil
}
