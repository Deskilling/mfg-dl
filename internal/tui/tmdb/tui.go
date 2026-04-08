package tmdb

import (
	"errors"
	"fmt"

	"mfg-dl/internal/search"
	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/tui/components"
	"mfg-dl/internal/util"

	"charm.land/log/v2"
)

func Tui() (err error) {
	if len(sites.Sites) == 0 {
		return errors.New("no service available")
	}

	tmdbResult, err := Search()
	if err != nil {
		return fmt.Errorf("failed to search tmdb: %w", err)
	}

	service, index, err := Score(tmdbResult)
	if err != nil {
		return fmt.Errorf("failed to calcualte scores: %w", err)
	}

	site := sites.Sites[index]
	result, err := sites.Sites[index].Search(tmdbResult.Score.Query[service])
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(result) == 0 {
		return errors.New("no search results")
	}

	season, err := components.Seasons(site, result[0])
	if err != nil {
		return fmt.Errorf("failed getting seasons: %w", err)
	}

	streams, err := components.Episodes(site, season)
	if err != nil || len(streams) == 0 {
		return fmt.Errorf("failed getting streams from episodes: %w", err)
	}

	err = components.DownloadMultiple(site, streams)
	if err != nil {
		return fmt.Errorf("failed downloading: %w", err)
	}

	log.Info("All Downloads finished")
	return nil
}

func Search() (result model.SearchResult, err error) {
	for {
		input, err := components.ReadString(components.Reader, "Search: ")
		if err != nil {
			return model.SearchResult{}, fmt.Errorf("failed userinput: %w", err)
		}

		tmdbResults, err := search.Search(input)
		if err != nil {
			continue
		}
		if len(tmdbResults) == 0 {
			continue
		}

		selected, err := selectFromList(tmdbResults)
		if err != nil {
			return model.SearchResult{}, fmt.Errorf("failed to select from list: %w", err)
		}

		var allServiceResults [][]model.SearchResult
		for _, site := range sites.Sites {
			results, err := site.Search(util.NormalizeString(selected.Name))
			if err != nil {
				log.Error("Search returned error", "service", site.Name(), "term", selected.Name, "err", err)
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
		return selected, nil
	}
}

func selectFromList(results []model.SearchResult) (result model.SearchResult, err error) {
	if len(results) == 0 {
		return model.SearchResult{}, errors.New("results are empty")
	}

	if len(results) == 1 {
		log.Infof("Selcted the only Result: %s", result.Name)
		return results[0], nil
	}

	for {
		for i, v := range results {
			fmt.Printf("[%v] %s %s\n", i+1, v.Name, v.ProductionYear)
		}

		input, err := components.ReadInt(components.Reader, "Enter: ")
		if err != nil {
			return model.SearchResult{}, fmt.Errorf("failed userinput: %w", err)
		}
		input--
		if input < 0 || input >= len(results) {
			util.ClearTerminal()
			log.Error("Invalid input", "min", 1, "max", len(results))
			continue
		}

		return results[input], nil
	}
}

func Score(result model.SearchResult) (service string, index int, err error) {
	if len(sites.Sites) == 1 {
		return sites.Sites[0].Name(), 0, nil
	}

	for {
		fmt.Printf("%s\n", result.Name)
		for i, v := range sites.Sites {
			// TODO add like color based on percentage
			// TODO filter out bad results
			fmt.Printf("[%v] %s %.2f%%\n", i+1, v.Name(), result.Score.Score[v.Name()]*100)
		}

		input, err := components.ReadInt(components.Reader, "Enter: ")
		if err != nil {
			return "", 0, fmt.Errorf("failed userinput: %w", err)
		}
		input--
		if input < 0 || input >= len(sites.Sites) {
			util.ClearTerminal()
			log.Error("Invalid input", "min", 1, "max", len(sites.Sites))
			continue
		}

		return sites.Sites[input].Name(), input, nil
	}
}
