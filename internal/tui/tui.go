package tui

import (
	"errors"
	"fmt"
	"log/slog"

	"mfg-dl/internal/search"
	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/tui/components"
)

func Tui() error {
	fmt.Println()
	slog.Debug("Starting Tui")

	tmdbResults, err := components.Search("Search", search.Search)
	if err != nil {
		return err
	}

	selectedResult, err := components.SingleSelect(tmdbResults, true, func(r model.SearchResult) string {
		if r.ProductionYear != "" {
			return fmt.Sprintf("%s (%s)", r.Name, r.ProductionYear)
		}

		return r.Name
	})
	if err != nil {
		return fmt.Errorf("failed selecting single search result: %w", err)
	}

	selectedResult, err = search.Match(selectedResult)
	if err != nil {
		return fmt.Errorf("failed calculating scores: %w", err)
	}

	var scoreOptions []string
	scoreByLabel := map[string]model.Score{}
	// currently it just appends in order of the loading process
	// i could maybe put the highest match on top but im not sure
	// typign for search should be enough
	// also generally i dont really know if the score system in that way
	// is good but its good enough for now
	// but no worries :)

	for _, module := range sites.Modules {
		candidates := selectedResult.Score[module.Name()]
		for i, candidate := range candidates {
			label := fmt.Sprintf("%s [%.1f%%] (%s)", candidate.Query, candidate.Score*100, module.Name())
			if i == len(candidates)-1 {
				label += "\n"
			}
			scoreOptions = append(scoreOptions, label)
			scoreByLabel[label] = candidate
		}
	}

	selectedSeries, err := components.SingleSelect(scoreOptions, true, func(r string) string {
		return r
	})
	if err != nil {
		return fmt.Errorf("failed selecting score: %w", err)
	}

	actualSelectedSeries, ok := scoreByLabel[selectedSeries]
	if !ok {
		return errors.New("selected result not found in options (???)")
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

	selectedSeasons, err := components.MultiSelect(seasons, false, func(r model.Season) string {
		return r.SeasonLabel
	})
	if err != nil {
		return fmt.Errorf("failed selecting seasons: %w", err)
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

		toDownloadEpisodes, err = components.MultiSelect(episodes, true, func(e model.Episode) string {
			return fmt.Sprintf("%s %s", e.EpisodeNum, e.EpisodeTitle)
		})
		if err != nil {
			return fmt.Errorf("failed selecting episodes: %w", err)
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

	slog.Debug("a", "a", toDownloadEpisodes)

	return nil
}
