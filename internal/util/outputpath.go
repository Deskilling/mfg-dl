package util

import (
	"mfg-dl/internal/core"
	"mfg-dl/internal/sites/model"

	"strings"
)

func BuildOutputPath(stream model.Stream) string {
	cfg := core.GetConfig()

	pattern := cfg.Extra.FilePattern

	replacer := strings.NewReplacer(
		"{location}", cfg.Location.Download,
		"{name}", stream.Name,
		"{season}", stream.SeasonNum,
		"{episode}", stream.EpisodeNum,
		"{language}", stream.Language,
		"{hoster},", stream.Hoster,
	)

	return replacer.Replace(pattern)
}
