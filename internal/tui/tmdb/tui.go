package tmdb

import (
	"fmt"

	"mfg-dl/internal/search"
	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/tui/components"
	"mfg-dl/internal/util"

	"charm.land/log/v2"
)

func Tui() {
	tmdbResult := Search()
	service, index := Score(tmdbResult)

	var site *model.Site
	site = &sites.Sites[index]
	result, _ := sites.Sites[index].Search(tmdbResult.Score.Query[service])
	if len(result) == 0 {
		log.Error("No search results")
		return
	}

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
				log.Error("Search returned error", "service", site.Service, "term", selected.Name, "err", err)
				continue
			}

			if len(results) == 0 {
				// Might to add switch to other tui, but then this whole thing wont make sense
				// this is why it got removed rn might add back in the future
				log.Debug("No match", "service", selected.Service)
			} else {
				log.Debug("Matches", "service", selected.Service)
				allServiceResults = append(allServiceResults, results)
			}
		}

		search.Match(&selected, allServiceResults)
		log.Debug("Got all Scores", "selected", selected.Score)
		return selected
	}
}

func selectFromList(results []model.SearchResult) model.SearchResult {
	if len(results) == 0 {
		log.Error("Search results are empty :(")
		return model.SearchResult{}
	}

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
			util.ClearTerminal()
			log.Error("Invalid input", "min", 1, "max", len(results))
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
		fmt.Printf("%s\n", result.Name)
		for i, v := range sites.Sites {
			// TODO add like color based on percentage
			// TODO filter out bad results
			fmt.Printf("[%v] %s %.2f%%\n", i+1, v.Service, result.Score.Score[v.Service]*100)
		}

		input := components.ReadInt(components.Reader, "Enter: ")
		input--
		if input < 0 || input >= len(sites.Sites) {
			util.ClearTerminal()
			log.Error("Invalid input", "min", 1, "max", len(sites.Sites))
			continue
		}

		return sites.Sites[input].Service, input
	}
}
