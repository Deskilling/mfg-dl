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

type Stream struct {
	Name                    string
	Href                    string
	SeasonNum               string
	EpisodeTitle            string
	EpisodeAlternativeTitle string
	EpisodeNum              string
	Hoster                  string
	Language                string
}

func GetStreams(episode model.Episode) (streams []model.Stream, err error) {
	pageURL := BaseURL + episode.Href
	log.Debug("leggo schmeggo", "pageURL", pageURL)
	unparsedStreams, err := request.Get(pageURL)
	if err != nil {
		err = fmt.Errorf("failed to GET Stream for %s: %w", episode.Href, err)
		log.Error(err)
		return nil, err
	}

	start := time.Now()
	parsedStreams, err := ParseStreams(unparsedStreams)
	if err != nil {
		err = fmt.Errorf("failed parsing Streams for %s: %w", episode.Href, err)
		log.Error(err)
		return nil, err
	}
	log.Debugf("time took for stream parsing: %v", time.Since(start))
	log.Debugf("parsed %d streams for %s", len(parsedStreams), episode.Href)

	if len(parsedStreams) == 0 {
		err = fmt.Errorf("%s not found", episode.Href)
		log.Error(err)
		return nil, err
	}

	for i := range parsedStreams {
		parsedStreams[i].Name = episode.Name
		parsedStreams[i].SeasonNum = episode.SeasonNum
		parsedStreams[i].EpisodeTitle = episode.EpisodeTitle
		parsedStreams[i].EpisodeAlternativeTitle = episode.EpisodeAlternativeTitle
		parsedStreams[i].EpisodeNum = episode.EpisodeNum
	}

	for _, v := range parsedStreams {
		stream := model.Stream{
			Name:                    v.Name,
			Href:                    v.Href,
			SeasonNum:               v.SeasonNum,
			EpisodeTitle:            v.EpisodeTitle,
			EpisodeAlternativeTitle: v.EpisodeAlternativeTitle,
			EpisodeNum:              v.EpisodeNum,
			Hoster:                  v.Hoster,
			Language:                v.Language,
		}

		streams = append(streams, stream)
	}

	return streams, nil
}

func ParseStreams(html string) (streams []Stream, err error) {
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
			Language: AniLanguages[lang],
		})
	})

	return streams, nil
}
