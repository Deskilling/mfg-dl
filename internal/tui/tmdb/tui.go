package tmdb

import (
	"fmt"
	"os"

	"mfg-dl/internal/search"
	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/tui/components"
	"mfg-dl/internal/tui/service"
	"mfg-dl/internal/util"
)

func Tui() {
	tmdbResult := Search()
	service, index := Score(tmdbResult)

	var site *model.Site
	site = &sites.Sites[index]
	result, _ := sites.Sites[index].Search(tmdbResult.Score.Query[service])
	season := components.Seasons(site, result[0])

	streams := components.Episodes(site, season)
	if len(streams) == 0 {
		return
	}

	site.DownloadMultiple(streams)
}

func Search() (result model.SearchResult) {
	for {
		input := components.ReadString(components.Reader, "Search: ")

		tmdbResults, err := search.Search(input)
		if err != nil {
			continue
		}
		if len(tmdbResults) == 0 {
			continue
		}

		selected := selectFromList(tmdbResults)

		var allServiceResults [][]model.SearchResult
		for _, site := range sites.Sites {
			results, err := site.Search(util.NormalizeString(selected.Name))
			if err != nil {
				continue
			}
			if len(results) == 0 {
				results, err = site.Search(util.ShortSearchTerm(util.NormalizeString(selected.Name)))
				if len(results) == 0 {
					input = components.ReadString(components.Reader, "y/n: ")
					if input == "y" {
						service.Tui()
					}
					os.Exit(0)
				}
			} else {
				allServiceResults = append(allServiceResults, results)
			}
		}

		search.Match(&selected, allServiceResults)
		return selected
	}
}

func selectFromList(results []model.SearchResult) model.SearchResult {
	if len(results) == 1 {
		return results[0]
	}
	for {
		for i, v := range results {
			fmt.Printf("[%v] %s %s\n", i+1, v.Name, v.ProductionYear)
		}
		input := components.ReadInt(components.Reader, "Enter: ")
		input--
		if input < 0 || input >= len(results) {
			continue
		}
		return results[input]
	}
}

func Score(result model.SearchResult) (service string, index int) {
	if len(sites.Sites) == 1 {
		return sites.Sites[0].Service, 0
	}

	for {
		for i, v := range sites.Sites {
			fmt.Printf("[%v] %s %.2f%%\n", i+1, v.Service, result.Score.Score[v.Service]*100)
		}
		input := components.ReadInt(components.Reader, "Enter: ")
		input--

		if input < 0 || input >= len(sites.Sites) {
			continue
		}

		return sites.Sites[input].Service, input
	}
}
