package aniworld

import (
	"time"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/sites/modules/voe"
	"mfg-dl/internal/util"

	"github.com/charmbracelet/log"
)

func Download(stream model.Stream) (err error) {
	url := BaseURL + stream.Href
	location, _ := request.Redirect(url)

	output := util.BuildOutputPath(stream)

	start := time.Now()
	voe.BaseDownload(location, output)
	log.Errorf("Took %s to download to %s", time.Since(start), location)

	return nil
}
