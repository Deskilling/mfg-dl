package serienstream

import (
	"strings"

	"mfg-dl/internal/sites/model"

	"charm.land/log/v2"
	"github.com/gocolly/colly/v2"
)

func (service *Serienstream) Streams(episode model.Episode) (streams []model.Stream, err error) {
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	c.OnHTML("#episode-links button.link-box", func(e *colly.HTMLElement) {
		href := e.Attr("data-play-url")
		if href == "" {
			return
		}

		hosterName := strings.TrimSpace(e.Attr("data-provider-name"))
		if hosterName == "" {
			return
		}

		language := strings.TrimSpace(e.Attr("data-language-label"))
		if language == "" {
			return
		}

		streams = append(streams, model.Stream{
			Service:                 Name,
			Name:                    episode.Name,
			Href:                    href,
			SeasonNum:               episode.SeasonNum,
			EpisodeTitle:            episode.EpisodeTitle,
			EpisodeAlternativeTitle: episode.EpisodeAlternativeTitle,
			EpisodeNum:              episode.EpisodeNum,
			Hoster:                  hosterName,
			Language:                language,
		})
	})

	err = c.Visit(BaseURL + episode.Href)
	if err != nil {
		log.Error(err)
	}
	c.Wait()

	return streams, nil
}
