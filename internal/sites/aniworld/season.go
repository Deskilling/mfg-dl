package aniworld

import (
	"fmt"
	"strings"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

func GetSeasons(result sites.SearchResult) ([]sites.Season, error) {
	seasons, err := request.Get(AniEndpoints["episodes"] + result.Href)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	parsedSeasons, err := ParseSeasons(seasons)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if len(parsedSeasons) == 0 {
		err = fmt.Errorf("%s not found", result.Href)
		log.Error(err)
		return nil, err
	}

	for i := range parsedSeasons {
		parsedSeasons[i].Name = result.Name
	}

	return parsedSeasons, nil
}

func ParseSeasons(html string) (seasons []sites.Season, err error) {
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

			seasons = append(seasons, sites.Season{
				Href:      href,
				Label:     label,
				SeasonNum: seasonNumber,
			})
		}
	})

	return seasons, nil
}
