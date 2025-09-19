package parserAniworld

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Stream struct {
	Href     string
	Hoster   string
	Url      string
	Language string
}

func Streams(html string) ([]Stream, error) {
	if html == "" {
		return nil, fmt.Errorf("not html parsed")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("could not create goquery document: %w", err)
	}

	var streams []Stream

	doc.Find("li[class*='episodeLink']").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a.watchEpisode")

		href, exists := link.Attr("href")
		if !exists {
			return // Skip this item if no link exists
		}

		hosterName := strings.TrimSpace(link.Find("h4").Text())

		if hosterName == "" {
			return
		}

		langKey, exists := s.Attr("data-lang-key")
		var lang string
		if exists {
			lang = strings.TrimSpace(langKey)
		}

		streams = append(streams, Stream{
			Href:     href,
			Hoster:   hosterName,
			Url:      href,
			Language: lang,
		})
	})

	return streams, nil
}
