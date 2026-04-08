package components

import (
	"mfg-dl/internal/sites/model"
	"time"

	"charm.land/log/v2"
)

func DownloadMultiple(service model.Site, streams []model.Stream) (err error) {
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
