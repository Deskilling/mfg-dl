package aniworld

import (
	"fmt"
	"strings"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"

	"github.com/PuerkitoBio/goquery"
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
		return nil, err
	}

	parsedEpisodes, err := ParseEpisodes(string(unparsedEpisodes))
	if err != nil {
		err = fmt.Errorf("failed parsing episodes %w", err)
		return nil, err
	}

	if len(parsedEpisodes) == 0 {
		return nil, err
	}

	for i := range parsedEpisodes {
		parsedEpisodes[i].Name = season.Name
		parsedEpisodes[i].SeasonNum = season.SeasonNum
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
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		err = fmt.Errorf("could not create goquery document: %w", err)
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
