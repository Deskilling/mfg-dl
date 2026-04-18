package tmdb

import (
	"errors"
	"fmt"

	"mfg-dl/internal/search"
	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/tui/components"
	"mfg-dl/internal/util"

	"charm.land/huh/v2"
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
		input, err := components.ReadString("Search: ")
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

	var values []huh.Option[int]
	for i, v := range results {
		values = append(values, huh.NewOption(v.Name, i))
	}

	var v int
	err = huh.NewSelect[int]().
		Title("Pick a Result").
		Options(values...).
		Value(&v).
		WithTheme(components.Theme).
		Run()
	if err != nil {
		log.Fatal(err)
	}

	return results[v], nil
}

func Score(result model.SearchResult) (service string, index int, err error) {
	if len(sites.Sites) == 1 {
		return sites.Sites[0].Name(), 0, nil
	}

	var values []huh.Option[int]
	for i, v := range sites.Sites {
		label := fmt.Sprintf("%s %.2f%%", v.Name(), result.Score.Score[v.Name()]*100)
		values = append(values, huh.NewOption(label, i))
	}

	var selected int
	err = huh.NewSelect[int]().
		Title(result.Name).
		Options(values...).
		Value(&selected).
		WithTheme(components.Theme).
		Run()
	if err != nil {
		return "", 0, fmt.Errorf("failed user input: %w", err)
	}

	return sites.Sites[selected].Name(), selected, nil
}
