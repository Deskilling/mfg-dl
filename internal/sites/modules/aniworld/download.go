package aniworld

import (
	"slices"
	"time"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/sites/modules/voe"
	"mfg-dl/internal/util"

	"github.com/charmbracelet/log"
)

func Download(stream model.Stream) (err error) {
	if !slices.Contains(Hoster, stream.Hoster) {
		return nil
	}

	url := BaseURL + stream.Href
	location, _ := request.Redirect(url)

	output := util.BuildOutputPath(stream)

	start := time.Now()

	switch stream.Hoster {
	case "VOE":
		voe.BaseDownload(location, output)
	default:
	}

	log.Infof("Time took for Download %s", time.Since(start))
	return nil
}

func DownloadMultiple(streams []model.Stream) (err error) {
	saft := time.Now()

	for _, v := range streams {
		Download(v)
	}

	log.Infof("Time took for all Downloads %s", time.Since(saft))
	return nil
}
