package aniworld

import (
	"fmt"
	"strings"

	"mfg-dl/request"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

type Episode struct {
	Href     string
	Title    string
	EngTitle string
}

func GetEpisodes(anime, season string) ([]Episode, error) {
	episodes, err := request.Get(AniEndpoints["episodes"] + anime + "/staffel-" + season)
	if err != nil {
		err = fmt.Errorf("failed to GET Episodes for %s: %w", anime, err)
		log.Error(err)
		return nil, err
	}

	parsedEpisodes, err := parseEpisodes(episodes)
	if err != nil {
		err = fmt.Errorf("failed parsing episodes for %s %s: %w", anime, season, err)
		log.Error(err)
		return nil, err
	}
	if len(parsedEpisodes) == 0 {
		err = fmt.Errorf("%s %s not found", anime, season)
		log.Error(err)
		return nil, err
	}

	return parsedEpisodes, nil
}

func parseEpisodes(html string) ([]Episode, error) {
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

	var episodes []Episode

	doc.Find(".seasonEpisodesList tbody tr").Each(func(i int, s *goquery.Selection) {
		episodeLink := s.Find("td.seasonEpisodeTitle a")
		href, exists := episodeLink.Attr("href")
		if !exists {
			return
		}

		title := strings.TrimSpace(episodeLink.Find("strong").Text())
		extra := strings.TrimSpace(episodeLink.Find("span").Text())

		episodes = append(episodes, Episode{
			Href:     href,
			Title:    title,
			EngTitle: extra,
		})
	})

	return episodes, nil
}
