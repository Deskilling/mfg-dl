package aniworld

import (
	"fmt"
	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/voe"
	"mfg-dl/internal/util"
	"mfg-dl/pkg/filesystem"
	"strings"

	"github.com/charmbracelet/log"
)

func Download(stream sites.Stream) error {
	anime, season, episode, err := util.ParseAnimePath(stream.Href)
	if err != nil {
		log.Error("failed parsing animehref", "err", err)
		return err
	}

	output := fmt.Sprintf(
		"%s%s/season-%v/%s-season%v-episode%v-%s.mp4",
		util.GetSettings().Location.Output,
		anime,
		season,
		anime,
		season,
		episode,
		AniLanguages[stream.Language],
	)

	if strings.Contains(stream.Href, "filme") {
		output = fmt.Sprintf(
			"%s%s/filme/%s-film%v-%s.mp4",
			util.GetSettings().Location.Output,
			anime,
			anime,
			episode,
			AniLanguages[stream.Language],
		)
	}

	if filesystem.ExistPath(output) {
		log.Info("already downloaded ... skipping", "season", season, "episode", episode)
		return nil
	}

	switch stream.Hoster {
	case "VOE":
		log.Info("Start Download", "season", season, "episode", episode)
		return voe.BaseDownload(BaseURL+stream.Href, output)

	case "Vidmoly", "Filemoon":
		return fmt.Errorf("hoster not implemented: %s", stream.Hoster)

	default:
		return fmt.Errorf("invalid hoster: %s", stream.Hoster)
	}
}
