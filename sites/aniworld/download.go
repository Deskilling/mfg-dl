package aniworld

import (
	"errors"
	"fmt"

	"mfg-dl/request"
	"mfg-dl/sites/voe"

	"github.com/charmbracelet/log"
)

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
