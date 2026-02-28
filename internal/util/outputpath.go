package util

import (
	"mfg-dl/internal/core"
	"mfg-dl/internal/sites/model"

	"strings"
)

func BuildOutputPath(stream model.Stream) string {
	cfg := core.GetConfig()

	pattern := cfg.Location.FilePattern

	replacer := strings.NewReplacer(
		"{location}", cfg.Location.Download,
		"{name}", stream.Name,
		"{season}", stream.SeasonNum,
		"{title}", stream.EpisodeTitle,
		"{title2}", stream.EpisodeAlternativeTitle,
		"{episode}", stream.EpisodeNum,
		"{language}", stream.Language,
		"{hoster},", stream.Hoster,
	)

	return replacer.Replace(pattern)
}
