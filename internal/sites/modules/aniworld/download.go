package aniworld

import (
	"slices"
	"time"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

	"charm.land/log/v2"
)

func (service *Aniworld) Download(stream model.Stream) (err error) {
	if !slices.Contains(Hoster, stream.Hoster) {
		return nil
	}

	url := BaseURL + stream.Href
	location, _ := request.Redirect(url)

	output := util.BuildOutputPath(stream)

	switch stream.Hoster {
	case "VOE":
		service.voe.BaseDownload(location, output)
	default:
	}

	return nil
}

func (service *Aniworld) DownloadMultiple(streams []model.Stream) (err error) {
	for _, v := range streams {
		start := time.Now()
		err = service.Download(v)
		if err != nil {
			log.Error("failed to download", "stream", streams)
		} else {
			log.Info("Downloaded", "title", v.Name, "season", v.SeasonNum, "episode", v.EpisodeNum, "time", time.Since(start))
		}
	}

	return nil
}
