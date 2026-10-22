package components

import (
	"fmt"

	"mfg-dl/internal/sites/model"

	"github.com/pterm/pterm"
)

type SearchFunc func(query string) ([]model.SearchResult, error)

func Search(prompt string, search SearchFunc) ([]model.SearchResult, error) {
	for {
		query, err := pterm.DefaultInteractiveTextInput.
			WithDefaultText(prompt).
			Show()
		if err != nil {
			return nil, fmt.Errorf("reading search input: %w", err)
		}

		results, err := search(query)
		if err != nil {
			return nil, fmt.Errorf("searching: %w", err)
		}

		if len(results) == 0 {
			pterm.Warning.Printfln("No results for %q", query)
			continue
		}

		return results, nil
	}
}
