package aniworld

import (
	"fmt"
	"strings"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

func GetEpisodes(season sites.Season) ([]sites.Episode, error) {
	url := AniEndpoints["episodes"] + season.Href

	episodes, err := request.Get(url)
	if err != nil {
		err = fmt.Errorf("failed to GET Episodes: %w", err)
		log.Error(err)
		return nil, err
	}

	parsedEpisodes, err := ParseEpisodes(episodes)
	if err != nil {
		err = fmt.Errorf("failed parsing episodes %w", err)
		log.Error(err)
		return nil, err
	}

	if len(parsedEpisodes) == 0 {
		log.Error(err)
		return nil, err
	}

	for i := range parsedEpisodes {
		parsedEpisodes[i].Name = season.Name
		parsedEpisodes[i].SeasonNum = season.SeasonNum

	}

	return parsedEpisodes, nil
}

func ParseEpisodes(html string) (episodes []sites.Episode, err error) {
	if html == "" {
		err := fmt.Errorf("not html parsed")
		log.Error(err)
		return nil, err
	}

	/*
		filesystem.WriteFile("./unparsed_episode.html", []byte(html))
		log.Fatal("quit")
	*/

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		err = fmt.Errorf("could not create goquery document: %w", err)
		log.Error(err)
		return nil, err
	}

	doc.Find(".seasonEpisodesList tbody tr").Each(func(i int, s *goquery.Selection) {
		episodeLink := s.Find("td.seasonEpisodeTitle a")
		href, exists := episodeLink.Attr("href")
		if !exists {
			return
		}
		episodeNum, exists := s.Find("meta[itemprop='episodeNumber']").Attr("content")
		if !exists {
			episodeNum = "0"
		}

		title := strings.TrimSpace(episodeLink.Find("strong").Text())
		//extra := strings.TrimSpace(episodeLink.Find("span").Text())

		episodes = append(episodes, sites.Episode{
			Href:       href,
			Title:      title,
			EpisodeNum: episodeNum,
			//EngTitle: extra,
		})
	})

	return episodes, nil
}
