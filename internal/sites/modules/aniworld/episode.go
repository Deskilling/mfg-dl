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

type Episode struct {
	Name                    string
	Href                    string
	SeasonNum               string
	EpisodeTitle            string
	EpisodeAlternativeTitle string
	EpisodeNum              string
}

func GetEpisodes(season model.Season) (episodes []model.Episode, err error) {
	url := BaseURL + season.Href

	unparsedEpisodes, err := request.Get(url)
	if err != nil {
		err = fmt.Errorf("failed to GET Episodes: %w", err)
		log.Error(err)
		return nil, err
	}

	start := time.Now()
	parsedEpisodes, err := ParseEpisodes(unparsedEpisodes)
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

	for _, v := range parsedEpisodes {
		episode := model.Episode{
			Name:                    v.Name,
			Href:                    v.Href,
			SeasonNum:               v.SeasonNum,
			EpisodeTitle:            v.EpisodeTitle,
			EpisodeAlternativeTitle: v.EpisodeAlternativeTitle,
			EpisodeNum:              v.EpisodeNum,
		}

		episodes = append(episodes, episode)
	}

	return episodes, nil
}

func ParseEpisodes(html string) (episodes []Episode, err error) {
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

	doc.Find(".seasonEpisodesList tbody tr").Each(func(i int, s *goquery.Selection) {
		episodeLink := s.Find("td.seasonEpisodeTitle a")
		href, exists := episodeLink.Attr("href")
		if !exists {
			return
		}
		episodeNum, exists := s.Find("meta[itemprop='episodeNumber']").Attr("content")
		if !exists {
			episodeNum = "00"
		} else {
			if len(episodeNum) == 1 {
				episodeNum = "0" + episodeNum
			}
		}

		title := strings.TrimSpace(episodeLink.Find("strong").Text())
		alternativeTitle := strings.TrimSpace(episodeLink.Find("span").Text())

		episodes = append(episodes, Episode{
			Href:                    href,
			EpisodeTitle:            title,
			EpisodeAlternativeTitle: alternativeTitle,
			EpisodeNum:              episodeNum,
		})
	})

	return episodes, nil
}
