package tui

import (
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

		return fmt.Sprintf("%s", r.Name)
	})

	selectedResult, err = search.Match(selectedResult)
	if err != nil {
		return fmt.Errorf("failed calculating scores: %w", err)
	}

	var scoreOptions []string

	// currently it just appends in order of the loading process
	// i could maybe put the highest match on top but im not sure
	// typign for search should be enough
	// also generally i dont really know if the score system in that way
	// is good but its good enough for now
	// but no worries :)
	for _, module := range sites.Modules {
		if len(selectedResult.Score[module.Name()]) == 0 {
			continue
		}
		for _, canidate := range selectedResult.Score[module.Name()] {
			scoreOptions = append(scoreOptions, fmt.Sprintf("%s [%0.1f%%] (%s)", canidate.Query, canidate.Score*100, module.Name()))
		}
		scoreOptions[len(scoreOptions)-1] += "\n"
	}

	components.SingleSelect(scoreOptions, true, func(r string) string {
		return r
	})

	slog.Debug("got result", "resultScores", selectedResult.Score)

	return nil
}
