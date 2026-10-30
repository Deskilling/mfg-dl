package serienstream

import (
	"fmt"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

	"slices"
)

func (service *Serienstream) Download(stream model.Stream) (err error) {
	if slices.Index(Hoster, stream.Hoster) < 0 {
		return nil
	}

	url := BaseURL + stream.Href
	location, err := request.Redirect(nil, url)
	if err != nil {
		return fmt.Errorf("failed to resolve redirect: %w", err)
	}

	output := util.BuildOutputPath(stream)

	switch stream.Hoster {
	case "VOE":
		return service.voe.BaseDownload(location, output)
	default:
	}

	return nil
}
