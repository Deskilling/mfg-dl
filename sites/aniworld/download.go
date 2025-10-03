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
	var mu sync.Mutex
	semaphore := make(chan struct{}, maxConcurrency)

	var failedSeason []string
	var failedEpisode []string

	for i, v := range season {
		episodes, err := GetEpisodes(anime, v)
		if err != nil {
			mu.Lock()
			failedSeason = append(failedSeason, strconv.Itoa(i))
			failedEpisode = append(failedEpisode, "0")
			mu.Unlock()
			continue
		}

		for j, k := range episodes {
			semaphore <- struct{}{}
			wg.Add(1)

			go func(i int, j int, k Episode) {
				defer wg.Done()
				defer func() { <-semaphore }()

				ep := strconv.Itoa(j + 1)
				err := Download(anime, v, ep, language, prefHost)
				if err != nil {
					mu.Lock()
					failedSeason = append(failedSeason, season[i])
					failedEpisode = append(failedEpisode, ep)
					mu.Unlock()
				}
			}(i, j, k)
		}
	}

	wg.Wait()

	mu.Lock()
	if len(failedSeason) > 0 && len(failedEpisode) > 0 {
		for i := range failedSeason {
			err := Download(anime, failedSeason[i], failedEpisode[i], language, prefHost)
			if err != nil {
				log.Error("Failed downloading", anime, "Season:", failedSeason[i], "Episode:", failedEpisode[i])
			}
		}
	}
	mu.Unlock()
}

func Download(anime, season, episode, language, prefHost string) error {
	/*
		Aniworld support 2 different versions
		its /filme/film-x
		or /staffel-0/episode-x
	*/

	if season == "0" {
		season = "filme"
	}

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
