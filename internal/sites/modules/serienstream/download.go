package serienstream

import (
	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

	"golang.org/x/exp/slices"
)

func (service *Serienstream) Download(stream model.Stream) (err error) {
	if !slices.Contains(Hoster, stream.Hoster) {
		return nil
	}

	url := BaseURL + stream.Href
	location, _ := request.Redirect(nil, url)

	output := util.BuildOutputPath(stream)

	switch stream.Hoster {
	case "VOE":
		service.voe.BaseDownload(location, output)
	default:
	}

	return nil
}
