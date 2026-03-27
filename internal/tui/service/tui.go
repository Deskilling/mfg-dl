package service

import (
	"fmt"

	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/tui/components"

	"charm.land/log/v2"
)

var site *model.Site

func Tui() {
	site = SelectModule()
	result := Search()
	season := components.Seasons(site, result)
	streams := components.Episodes(site, season)

	if len(streams) == 0 {
		log.Error("no streams selected")
		return
	}

	site.DownloadMultiple(streams)
}

func SelectModule() *model.Site {
	if len(sites.Sites) == 1 {
		log.Infof("Selected %s", sites.Sites[0].Service)
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
			log.Error("Invalid Input, try again")
			continue
		}

		return &sites.Sites[v]
	}
}

func Search() model.SearchResult {
	for {
		search := components.ReadString(components.Reader, "Search: ")

		results, err := site.Search(search)
		if err != nil {
			log.Error("search failed", "err", err)
			continue
		}

		if len(results) == 0 {
			log.Error("Not Found, try again")
			continue
		}

		if len(results) == 1 {
			log.Infof("Selected %s", results[0].Name)
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
				log.Error("Invalid Input, try again")
				continue
			}

			return results[v]
		}
	}
}
