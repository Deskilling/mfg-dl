package aniworld

import (
	"fmt"
	"strings"

	"mfg-dl/request"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

type Season struct {
	Href         string
	Label        string
	SeasonNumber string
}

func GetSeasons(anime string) ([]Season, error) {
	seasons, err := request.Get(AniEndpoints["episodes"] + anime)
	if err != nil {
		err = fmt.Errorf("failed to GET Seasons for %s with error %w", anime, err)
		log.Error(err)
		return nil, err
	}

	parsedSeasons, err := parseSeasons(seasons)
	if err != nil {
		err = fmt.Errorf("failed parsing seasons for %s: %w", anime, err)
		log.Error(err)
		return nil, err
	}
	if len(parsedSeasons) == 0 {
		err = fmt.Errorf("%s not found", anime)
		log.Error(err)
		return nil, err
	}

	return parsedSeasons, nil
}

func parseSeasons(html string) ([]Season, error) {
	if html == "" {
		err := fmt.Errorf("not html parsed")
		log.Error(err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		err = fmt.Errorf("could not create goquery document: %w", err)
		log.Error(err)
		return nil, err
	}

	var seasons []Season
	doc.Find(".hosterSiteDirectNav ul a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		} else if strings.Contains(href, "/staffel-") && !strings.Contains(href, "/episode-") || strings.Contains(href, "/filme") {
			label, exists := s.Attr("title")
			if !exists {
				return
			}

			seasonNumber := strings.TrimSpace(strings.TrimPrefix(label, "Staffel "))

			seasons = append(seasons, Season{
				Href:         href,
				Label:        label,
				SeasonNumber: seasonNumber,
			})
		}
	})

	return seasons, nil
}
