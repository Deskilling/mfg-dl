package aniworld

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"mfg-dl/request"
	"mfg-dl/sites/voe"

	"github.com/charmbracelet/log"
)

var maxConcurrency = 4

func DownloadSeason(anime, language, prefHost string, season []string) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrency)

	var failedSeason []string
	var failedEpisode []string

	for i, v := range season {
		episodes, err := GetEpisodes(anime, v)
		if err != nil {
			failedSeason = append(failedSeason, strconv.Itoa(i))
			failedEpisode = append(failedEpisode, "0")
			continue
		}

		for j, k := range episodes {
			semaphore <- struct{}{}
			wg.Add(1)

			go func(j int, k Episode) {
				defer wg.Done()
				defer func() { <-semaphore }()

				ep := strconv.Itoa(j)
				err := Download(anime, v, ep, language, prefHost)
				if err != nil {
					failedSeason = append(failedSeason, season[i])
					failedEpisode = append(failedEpisode, strconv.Itoa(j))

				}
			}(j, k)
		}
	}

	wg.Wait()

	if len(failedSeason) >= 1 && len(failedEpisode) >= 1 {
		for i := range failedSeason {
			Download(anime, failedSeason[i], failedEpisode[i], language, prefHost)
		}
	}
}

func Download(anime, season, episode, language, prefHost string) error {
	output := fmt.Sprintf("./downloads/%s/season-%s/%s-%s-episode-%s.mp4", anime, season, anime, AniLanguages[language], episode)

	streams, err := GetStreams(anime, season, episode)
	if err != nil {
		err = fmt.Errorf("failed to get streams: %w", err)
		log.Error(err)
		return err
	}

	var found bool = false
	var host string = "VOE"
	var index int = 0
	for i, v := range streams {
		if v.Hoster == prefHost && v.Language == language {
			found = true
			index = i
			host = prefHost
			break
		}
	}

	// remove lang check
	if !found {
		for i, v := range streams {
			if v.Hoster == prefHost {
				found = true
				index = i
				host = prefHost
				break
			}
		}
	}

	redirectUrl, err := request.Redirect(BaseURL + streams[index].Href)
	if err != nil {
		log.Error(err)
		return err
	}

	switch streams[index].Hoster {
	case host:
		voe.BaseDownload(redirectUrl, output)
	default:
		return errors.New("failed finding host")
	}

	return nil
}
