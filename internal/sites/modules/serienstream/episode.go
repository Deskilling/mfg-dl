package serienstream

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

func (service *Serienstream) Episodes(season model.Season) (episodes []model.Episode, err error) {
	url := BaseURL + season.Href

	unparsedEpisodes, err := request.Get(service.client, url)
	if err != nil {
		return nil, fmt.Errorf("failed to GET Episodes: %w", err)
	}

	parsedEpisodes, err := ParseEpisodes(string(unparsedEpisodes))
	if err != nil {
		return nil, fmt.Errorf("failed parsing episodes %w", err)
	}

	if len(parsedEpisodes) == 0 {
		return nil, fmt.Errorf("no episodes found for %s", season.Href)
	}

	for i := range parsedEpisodes {
		parsedEpisodes[i].Name = season.Name
		parsedEpisodes[i].SeasonNum = season.SeasonNum
	}

	for _, v := range parsedEpisodes {
		if v.EpisodeTitle == "" {
			v.EpisodeTitle = v.EpisodeAlternativeTitle
		}

		episodes = append(episodes, model.Episode{
			Service:                 Name,
			Name:                    v.Name,
			Href:                    v.Href,
			SeasonNum:               v.SeasonNum,
			EpisodeTitle:            v.EpisodeTitle,
			EpisodeAlternativeTitle: v.EpisodeAlternativeTitle,
			EpisodeNum:              v.EpisodeNum,
		})
	}

	return episodes, nil
}

func ParseEpisodes(html string) (episodes []Episode, err error) {
	if html == "" {
		return nil, fmt.Errorf("empty html")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("could not create goquery document: %w", err)
	}

	doc.Find(".episode-table tbody tr.episode-row").Each(func(i int, s *goquery.Selection) {
		onclick, exists := s.Attr("onclick")
		if !exists {
			return
		}

		href := strings.TrimPrefix(onclick, "window.location='")
		href = strings.TrimSuffix(href, "'")

		episodeNum := strings.TrimSpace(s.Find(".episode-number-cell").Text())
		if episodeNum == "" {
			episodeNum = "00"
		} else if len(episodeNum) == 1 {
			episodeNum = "0" + episodeNum
		}

		title := strings.TrimSpace(s.Find(".episode-title-ger").AttrOr("title", ""))
		alternativeTitle := strings.TrimSpace(s.Find(".episode-title-eng").AttrOr("title", ""))

		episodes = append(episodes, Episode{
			Href:                    href,
			EpisodeTitle:            title,
			EpisodeAlternativeTitle: alternativeTitle,
			EpisodeNum:              episodeNum,
		})
	})

	return episodes, nil
}
