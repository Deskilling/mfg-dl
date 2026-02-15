package aniworld

import (
	"fmt"
	"strings"
	"time"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

func GetEpisodes(season model.Season) ([]model.Episode, error) {
	url := BaseURL + season.Href

	episodes, err := request.Get(url)
	if err != nil {
		err = fmt.Errorf("failed to GET Episodes: %w", err)
		log.Error(err)
		return nil, err
	}

	start := time.Now()
	parsedEpisodes, err := ParseEpisodes(episodes)
	if err != nil {
		err = fmt.Errorf("failed parsing episodes %w", err)
		log.Error(err)
		return nil, err
	}
	log.Debugf("time took for episode parsing: %v", time.Since(start))

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

func ParseEpisodes(html string) (episodes []model.Episode, err error) {
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

		episodes = append(episodes, model.Episode{
			Href:       href,
			Title:      title,
			EpisodeNum: episodeNum,
			//EngTitle: extra,
		})
	})

	return episodes, nil
}
