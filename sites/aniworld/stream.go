package aniworld

import (
	"fmt"
	"strings"

	"mfg-dl/request"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

type Stream struct {
	Href     string
	Hoster   string
	Url      string
	Language string
}

func GetStreams(anime, season, episode string) ([]Stream, error) {
	pageURL := AniEndpoints["episodes"] + anime + "/staffel-" + season + "/episode-" + episode
	log.Debug(pageURL)
	streams, err := request.Get(pageURL)
	if err != nil {
		err = fmt.Errorf("failed to GET Stream for %s %s %s: %w", anime, season, episode, err)
		log.Error(err)
		return nil, err
	}

	parsedStreams, err := parseStreams(streams)
	if err != nil {
		err = fmt.Errorf("failed parsing Streams for %s %s %s: %w", anime, season, episode, err)
		log.Error(err)
		return nil, err
	}
	log.Debugf("parsed %d streams for %s %s %s", len(parsedStreams), anime, season, episode)
	if len(parsedStreams) == 0 {
		err = fmt.Errorf("%s not found", anime)
		log.Error(err)
		return nil, err
	}

	return parsedStreams, nil
}

func parseStreams(html string) ([]Stream, error) {
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

	var streams []Stream

	doc.Find("li[class*='episodeLink']").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a.watchEpisode")

		href, exists := link.Attr("href")
		if !exists {
			return
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
