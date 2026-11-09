package tui

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"mfg-dl/internal/search"
	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/tui/components"

	"github.com/pterm/pterm"
)

func breadcrumbs(parts ...string) {
	fmt.Print("\033[H\033[2J")
	pterm.FgGray.Println(strings.Join(parts, " › "))
	fmt.Println()
}

func trail(name, service string, parts ...string) {
	breadcrumbs(append([]string{"Search", name, service}, parts...)...)
}

func Tui() error {
	fmt.Println()
	slog.Debug("Starting Tui")

	breadcrumbs("Search")
	tmdbResults, err := components.Search("Search", search.Search)
	if err != nil {
		return err
	}

	var selectedResult model.SearchResult
	var actualSelectedSeries model.Score

search:
	for {
		if len(tmdbResults) == 1 {
			selectedResult = tmdbResults[0]
		} else {
			breadcrumbs("Search", "Results")
			selectedResult, err = components.SingleSelect(tmdbResults, true, func(r model.SearchResult) string {
				if r.ProductionYear != "" {
					return fmt.Sprintf("%s (%s)", r.Name, r.ProductionYear)
				}

				return r.Name
			})
			if err != nil {
				return fmt.Errorf("failed selecting single search result: %w", err)
			}
		}

		selectedResult, err = search.Match(selectedResult)
		if err != nil {
			return fmt.Errorf("failed calculating scores: %w", err)
		}

		for {
			var overview []model.Score
			for _, m := range sites.Modules {
				candidates := selectedResult.Score[m.Name()]
				if len(candidates) > 0 {
					overview = append(overview, candidates[0])
				}
			}
			overview = append(overview, model.Score{Query: "← Back"})

			breadcrumbs("Search", selectedResult.Name)
			label := func(c model.Score) string {
				if c.Service == "" {
					return c.Query
				}
				m := sites.ModuleByName[c.Service]
				more := ""
				if n := len(selectedResult.Score[c.Service]); n > 1 {
					more = fmt.Sprintf(" +%d more", n-1)
				}
				return fmt.Sprintf("%s %s [%.1f%%]%s (Dubbed: %s | Subbed: %s)", pterm.FgGreen.Sprintf("%-13s", c.Service), c.Query, c.Score*100, more, m.Languages().Dubbed, m.Languages().Subbed)
			}

			top, err := components.SingleSelect(overview, true, label)
			if err != nil {
				return fmt.Errorf("failed selecting service: %w", err)
			}

			if top.Service == "" {
				if len(tmdbResults) == 1 {
					return nil
				}
				continue search
			}

			candidates := selectedResult.Score[top.Service]
			if len(candidates) == 1 {
				actualSelectedSeries = candidates[0]
				break search
			}

			breadcrumbs("Search", selectedResult.Name, top.Service)
			options := append(slices.Clone(candidates), model.Score{Query: "← Back"})
			picked, err := components.SingleSelect(options, true, func(c model.Score) string {
				if c.Service == "" {
					return c.Query
				}
				return fmt.Sprintf("%s [%.1f%%]", c.Query, c.Score*100)
			})
			if err != nil {
				return fmt.Errorf("failed selecting score: %w", err)
			}

			if picked.Service == "" {
				continue
			}

			actualSelectedSeries = picked
			break search
		}
	}

	module := sites.ModuleByName[actualSelectedSeries.Service]

	seasons, err := module.Seasons(model.SearchResult{
		Service: module.Name(),
		Name:    selectedResult.Name,
		Href:    actualSelectedSeries.Href,

		// only the three above are essential
		// these are just extra
		Cover:          selectedResult.Cover,
		CoverPath:      selectedResult.CoverPath,
		ProductionYear: selectedResult.ProductionYear,
	})
	if err != nil {
		return fmt.Errorf("failed getting seasons: %w", err)
	}

	var selectedSeasons []model.Season
	if len(seasons) > 1 {
		trail(selectedResult.Name, module.Name(), "Seasons")
		selectedSeasons, err = components.MultiSelect(seasons, false, func(r model.Season) string {
			return r.SeasonLabel
		})
		if err != nil {
			return fmt.Errorf("failed selecting seasons: %w", err)
		}
	}
	if len(selectedSeasons) == 0 {
		return errors.New("invalid seasons selection (expected atleast one)")
	}

	var toDownloadEpisodes []model.Episode
	if len(selectedSeasons) == 1 {
		episodes, err := module.Episodes(selectedSeasons[0])
		if err != nil {
			return fmt.Errorf("failed getting episodes for %s: %w", module.Name(), err)
		}

		var toDownloadEpisodes []model.Episode
		if len(episodes) > 1 {
			trail(selectedResult.Name, module.Name(), "Episodes")
			toDownloadEpisodes, err = components.MultiSelect(episodes, true, func(e model.Episode) string {
				return fmt.Sprintf("%s %s", e.EpisodeNum, e.EpisodeTitle)
			})
			if err != nil {
				return fmt.Errorf("failed selecting episodes: %w", err)
			}
		}
		if len(toDownloadEpisodes) == 0 {
			return errors.New("invalid episodes selection (expected atleast one)")
		}
	} else {
		for _, season := range selectedSeasons {
			episodes, err := module.Episodes(season)
			if err != nil {
				return fmt.Errorf("failed getting episodes for %s: %w", module.Name(), err)
			}

			toDownloadEpisodes = append(toDownloadEpisodes, episodes...)
		}
	}

	// idk how this could happen
	if len(toDownloadEpisodes) < 1 {
		return fmt.Errorf("yolo das war wohl richtig doof von dir :(")
	}

	var toDownloadStreams []model.Stream

	var selectedLanguages []string
	var preferredHosters []string

	for _, episode := range toDownloadEpisodes {
		stream, err := module.Streams(episode)
		if err != nil {
			slog.Error("failed getting streams for", "episode", episode.EpisodeNum, "err", err)
			slog.Warn("skipping")
			continue
		}

		seen := make(map[string]struct{})
		var languages []string
		for _, v := range stream {
			if _, ok := seen[v.Language]; !ok {
				seen[v.Language] = struct{}{}
				languages = append(languages, v.Language)
			}
		}

		if len(languages) == 0 {
			return errors.New("gg")
		}

		if selectedLanguages == nil {
			if len(languages) == 1 {
				selectedLanguages = languages
			} else {
				trail(selectedResult.Name, module.Name(), "Language")
				selectedLanguages, err = components.MultiSelect(languages, false, func(l string) string {
					return l
				})
				if err != nil {
					return fmt.Errorf("failed language select: %w", err)
				}
				if len(selectedLanguages) == 0 {
					return errors.New("no language selected")
				}
			}
		}

		if preferredHosters == nil {
			seenHoster := make(map[string]struct{})
			var hosters []string
			for _, v := range stream {
				if !slices.Contains(selectedLanguages, v.Language) {
					continue
				}
				if _, ok := seenHoster[v.Hoster]; !ok {
					seenHoster[v.Hoster] = struct{}{}
					hosters = append(hosters, v.Hoster)
				}
			}

			if len(hosters) == 1 {
				preferredHosters = hosters
			} else if len(hosters) > 1 {
				trail(selectedResult.Name, module.Name(), "Hoster")
				preferredHosters, err = components.MultiSelect(hosters, false, func(h string) string {
					return h
				})
				if err != nil {
					return fmt.Errorf("failed hoster select: %w", err)
				}
			}
		}

		var chosen *model.Stream
		bestLang, bestRank := 0, 0
		for i := range stream {
			langRank := slices.Index(selectedLanguages, stream[i].Language)
			if langRank == -1 {
				continue
			}
			rank := slices.Index(preferredHosters, stream[i].Hoster)
			if rank == -1 {
				rank = len(preferredHosters)
			}
			if chosen == nil || langRank < bestLang || (langRank == bestLang && rank < bestRank) {
				chosen = &stream[i]
				bestLang = langRank
				bestRank = rank
			}
		}

		if chosen == nil {
			slog.Warn("skip ggs")
			continue
		}
		toDownloadStreams = append(toDownloadStreams, *chosen)
	}

	for _, stream := range toDownloadStreams {
		err = module.Download(stream)
		if err != nil {
			slog.Error("Downloading", "stream", stream, "err", err)
			continue
		}

		pterm.Success.Println(
			fmt.Sprintf("%s %s %s - %s [%s - %s]",
				stream.Name,
				stream.SeasonNum,
				stream.EpisodeNum,
				stream.EpisodeTitle,
				stream.Language,
				stream.Hoster,
			),
		)
	}

	return nil
}
