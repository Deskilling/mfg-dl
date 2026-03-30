package aniworld

import (
	"slices"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/sites/modules/voe"
	"mfg-dl/internal/util"

	"charm.land/log/v2"
)

func Download(stream model.Stream) (err error) {
	if !slices.Contains(Hoster, stream.Hoster) {
		return nil
	}

	url := BaseURL + stream.Href
	location, _ := request.Redirect(url)

	output := util.BuildOutputPath(stream)

	switch stream.Hoster {
	case "VOE":
		voe.BaseDownload(location, output)
	default:
	}

	return nil
}

func DownloadMultiple(streams []model.Stream) (err error) {
	for _, v := range streams {
		err = Download(v)
		if err != nil {
			log.Error("failed to download", "stream", streams)
		}
	}

	return nil
}
