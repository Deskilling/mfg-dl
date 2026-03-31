package service

import (
	"fmt"

	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/tui/components"
	"mfg-dl/internal/util"

	"charm.land/log/v2"
)

func Tui() (err error) {
	var site *model.Site
	site, err = SelectModule()
	if err != nil {

	}
	result, err := Search(site)
	if err != nil {

	}
	season, err := components.Seasons(site, result)
	if err != nil {
	}

	streams, err := components.Episodes(site, season)
	if err != nil {
	}

	if len(streams) == 0 {
		return
	}

	err = site.DownloadMultiple(streams)
	if err != nil {

	}

	return nil
}

func SelectModule() (site *model.Site, err error) {
	if len(sites.Sites) == 1 {
		return &sites.Sites[0], nil
	}

	for i := range sites.Sites {
		u := i + 1
		fmt.Printf("[%v] %s\n", u, sites.Sites[i].Service)
	}

	for {
		v, err := components.ReadInt(components.Reader, "Enter: ")
		if err != nil {
			return &model.Site{}, fmt.Errorf("failed userinput: %w", err)
		}
		v--

		if v < 0 || v >= len(sites.Sites) {
			continue
		}

		return &sites.Sites[v], nil
	}
}

func Search(site *model.Site) (result model.SearchResult, err error) {
	for {
		search, err := components.ReadString(components.Reader, "Search: ")
		if err != nil {
			return model.SearchResult{}, fmt.Errorf("failed userinput: %w", err)
		}

		results, err := site.Search(search)
		if err != nil {
			log.Error("Failed to search", "service", site.Service, "term", search)
			continue
		}

		if len(results) == 0 {
			continue
		}

		if len(results) == 1 {
			return results[0], nil
		}

		for {
			for i := range results {
				u := i + 1
				fmt.Printf("[%v] %s\n", u, results[i].Name)
			}

			v, err := components.ReadInt(components.Reader, "Enter: ")
			if err != nil {
				return model.SearchResult{}, fmt.Errorf("failed userinput: %w", err)
			}
			v--

			if v < 0 || v >= len(results) {
				util.ClearTerminal()
				log.Error("Invalid selection")
				continue
			}

			return results[v], nil
		}
	}
}
