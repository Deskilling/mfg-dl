package components

import (
	"errors"
	"fmt"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

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
		return seasons[0], errors.New("no season found")
	}

	for {
		u := 1
		if seasons[0].SeasonLabel == "Alle Filme" {
			u--
		}

		for i := range seasons {
			fmt.Printf("[%v] %s\n", u, seasons[i].SeasonLabel)
			u++
		}

		v, err := ReadInt(Reader, "Select: ")
		if err != nil {
			return model.Season{}, fmt.Errorf("failed userinput: %s", err)
		}
		if seasons[0].SeasonLabel != "Alle Filme" {
			v--
		}

		if v < 0 || v >= len(seasons) {
			util.ClearTerminal()
			log.Error("Invalid selection")
			continue
		}

		return seasons[v], nil
	}
}
