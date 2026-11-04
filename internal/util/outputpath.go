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
		"{name}", clean(stream.Name),
		"{season}", stream.SeasonNum,
		"{title}", clean(stream.EpisodeTitle),
		"{title2}", clean(stream.EpisodeAlternativeTitle),
		"{episode}", stream.EpisodeNum,
		"{language}", stream.Language,
		"{hoster}", stream.Hoster,
	)

	return replacer.Replace(pattern)
}

func clean(s string) string {
	// TODO check on windows what should also be replaced
	return strings.ReplaceAll(s,
		// spaces next to ∕ because it looks better
		"/", " ∕ ",
	)
}
