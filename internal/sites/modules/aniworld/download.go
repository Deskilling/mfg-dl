package aniworld

import (
	"fmt"
	"mfg-dl/internal/core"
	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/sites/modules/voe"
)

func Download(stream model.Stream) (err error) {
	url := BaseURL + stream.Href
	location, _ := request.Redirect(url)

	output := fmt.Sprintf("%s/%s/%s-Season%s-Episode%s.mp4", core.GetConfig().Location.Download, stream.Name, stream.Name, stream.SeasonNum, stream.EpisodeNum)
	voe.BaseDownload(location, output)

	return nil
}
