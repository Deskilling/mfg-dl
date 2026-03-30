package components

import (
	"fmt"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

	"charm.land/log/v2"
)

func Seasons(site *model.Site, result model.SearchResult) model.Season {
	seasons, err := site.Seasons(result)
	if err != nil {
		return model.Season{}
	}

	for {
		if len(seasons) == 1 {
			return seasons[0]
		}

		u := 1
		if seasons[0].SeasonLabel == "Alle Filme" {
			u--
		}

		for i := range seasons {
			fmt.Printf("[%v] %s\n", u, seasons[i].SeasonLabel)
			u++
		}

		v := ReadInt(Reader, "Select: ")
		if seasons[0].SeasonLabel != "Alle Filme" {
			v--
		}

		if v < 0 || v >= len(seasons) {
			util.ClearTerminal()
			log.Error("Invalid selection")
			continue
		}

		return seasons[v]
	}
}
