package service

import (
	"fmt"

	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/tui/components"
)

func Tui() {
	var site *model.Site
	site = SelectModule()
	result := Search(site)
	season := components.Seasons(site, result)
	streams := components.Episodes(site, season)

	if len(streams) == 0 {
		return
	}

	site.DownloadMultiple(streams)
}

func SelectModule() *model.Site {
	if len(sites.Sites) == 1 {
		return &sites.Sites[0]
	}

	for i := range sites.Sites {
		u := i + 1
		fmt.Printf("[%v] %s\n", u, sites.Sites[i].Service)
	}

	for {
		v := components.ReadInt(components.Reader, "Enter: ")
		v--

		if v < 0 || v >= len(sites.Sites) {
			continue
		}

		return &sites.Sites[v]
	}
}

func Search(site *model.Site) model.SearchResult {
	for {
		search := components.ReadString(components.Reader, "Search: ")

		results, err := site.Search(search)
		if err != nil {
			continue
		}

		if len(results) == 0 {
			continue
		}

		if len(results) == 1 {
			return results[0]
		}

		for {
			for i := range results {
				u := i + 1
				fmt.Printf("[%v] %s\n", u, results[i].Name)
			}

			v := components.ReadInt(components.Reader, "Enter: ")
			v--

			if v < 0 || v >= len(results) {
				continue
			}

			return results[v]
		}
	}
}
