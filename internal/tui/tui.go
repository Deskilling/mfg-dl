package tui

import (
	"fmt"
	"log/slog"

	"mfg-dl/internal/search"
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

	_, err = components.SingleSelect(tmdbResults, func(r model.SearchResult) string {
		if r.ProductionYear != "" {
			return fmt.Sprintf("%s (%s)", r.Name, r.ProductionYear)
		}

		return fmt.Sprintf("%s", r.Name)
	})

	return nil
}
