package serienstream

import (
	"fmt"
	"strings"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"

	"github.com/PuerkitoBio/goquery"
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

func (service *Serienstream) Streams(episode model.Episode) (streams []model.Stream, err error) {
	pageURL := BaseURL + episode.Href
	unparsedStreams, err := request.Get(service.client, pageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to GET Stream for %s: %w", episode.Href, err)
	}

	parsedStreams, err := ParseStreams(string(unparsedStreams))
	if err != nil {
		return nil, fmt.Errorf("failed parsing Streams for %s: %w", episode.Href, err)
	}

	if len(parsedStreams) == 0 {
		return nil, fmt.Errorf("%s not found", episode.Href)
	}

	for i := range parsedStreams {
		parsedStreams[i].Name = episode.Name
		parsedStreams[i].SeasonNum = episode.SeasonNum
		parsedStreams[i].EpisodeTitle = episode.EpisodeTitle
		parsedStreams[i].EpisodeAlternativeTitle = episode.EpisodeAlternativeTitle
		parsedStreams[i].EpisodeNum = episode.EpisodeNum
	}

	for _, v := range parsedStreams {
		streams = append(streams, model.Stream{
			Name:                    v.Name,
			Href:                    v.Href,
			SeasonNum:               v.SeasonNum,
			EpisodeTitle:            v.EpisodeTitle,
			EpisodeAlternativeTitle: v.EpisodeAlternativeTitle,
			EpisodeNum:              v.EpisodeNum,
			Hoster:                  v.Hoster,
			Language:                v.Language,
		})
	}

	return streams, nil
}

func ParseStreams(html string) (streams []Stream, err error) {
	if html == "" {
		return nil, fmt.Errorf("empty html")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("could not create goquery document: %w", err)
	}

	doc.Find("#episode-links button.link-box").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("data-play-url")
		if !exists || href == "" {
			return
		}

		hosterName := strings.TrimSpace(s.AttrOr("data-provider-name", ""))

		if hosterName == "" {
			return
		}

		language := strings.TrimSpace(s.AttrOr("data-language-label", ""))

		streams = append(streams, Stream{
			Href:     href,
			Hoster:   hosterName,
			Language: language,
		})
	})

	return streams, nil
}
