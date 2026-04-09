package components

import (
	"fmt"
	"mfg-dl/internal/sites/model"

	"charm.land/huh/v2"
	"charm.land/log/v2"
)

func Seasons(site model.Site, result model.SearchResult) (season model.Season, err error) {
	seasons, err := site.Seasons(result)
	if err != nil {
		return model.Season{}, fmt.Errorf("failed getting seasons: %w", err)
	}

	if len(seasons) == 0 {
		log.Error("No season found", "site", site.Name(), "href", result.Href)
		return model.Season{}, nil
	}

	if len(seasons) == 1 {
		log.Info("Selected the only season available", "seasonNum", seasons[0].SeasonNum)
		return seasons[0], nil
	}

	var values []huh.Option[int]
	for i, v := range seasons {
		values = append(values, huh.NewOption(v.SeasonLabel, i))
	}

	var v int
	err = huh.NewSelect[int]().
		Title("Select Season").
		Options(values...).
		Value(&v).
		WithTheme(Theme).
		Run()
	if err != nil {
		return model.Season{}, fmt.Errorf("failed userinput: %w", err)
	}

	return seasons[v], nil
}
