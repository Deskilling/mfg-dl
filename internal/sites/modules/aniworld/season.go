package aniworld

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

func (service *Aniworld) Seasons(result model.SearchResult) (seasons []model.Season, err error) {
	unparsedSeasons, err := request.Get(service.client, BaseURL+result.Href)
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
			Service:     Name,
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
		err := fmt.Errorf("empty html")
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("could not create goquery document: %w", err)
	}

	doc.Find(".hosterSiteDirectNav ul:first-child a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		isFilme := strings.Contains(href, "/filme")
		isSeason := strings.Contains(href, "/staffel-")

		if !isFilme && !isSeason {
			return
		}

		label, exists := s.Attr("title")
		if !exists {
			return
		}

		var seasonNumber string
		if isFilme {
			seasonNumber = "00"
		} else {
			seasonNumber = strings.TrimSpace(strings.TrimPrefix(label, "Staffel "))
			if len(seasonNumber) == 1 {
				seasonNumber = "0" + seasonNumber
			}
		}

		seasons = append(seasons, Season{
			Href:        href,
			SeasonLabel: label,
			SeasonNum:   seasonNumber,
		})
	})

	return seasons, nil
}
