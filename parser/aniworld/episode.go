package parserAniworld

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

type Episode struct {
	Href     string
	Title    string
	EngTitle string
}

func Episodes(html string) ([]Episode, error) {
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
