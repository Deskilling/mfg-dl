package service

import (
	"errors"
	"fmt"

	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/tui/components"

	"charm.land/huh/v2"
	"charm.land/log/v2"
)

func Tui() (err error) {
	if len(sites.Sites) == 0 {
		return errors.New("no service available")
	}

	site, err := SelectModule()
	if err != nil {
		return errors.New("failed to select service")
	}
	result, err := Search(site)
	if err != nil {
		return fmt.Errorf("failed to search: %w", err)
	}
	season, err := components.Seasons(site, result)
	if err != nil {
		return fmt.Errorf("failed to get seasons: %w", err)
	}

	streams, err := components.Episodes(site, season)
	if err != nil {
		return fmt.Errorf("failed to get episodes: %w", err)
	}

	if len(streams) == 0 {
		return errors.New("no streams found")
	}

	err = components.DownloadMultiple(site, streams)
	if err != nil {
		// TODO add to downloadmultiple probabbly a check if everything failed or smth
		return fmt.Errorf("failed to donwload all streams")
	}

	return nil
}

func SelectModule() (site model.Site, err error) {
	if len(sites.Sites) == 1 {
		return sites.Sites[0], nil
	}

	for i, v := range sites.Sites {
		u := i + 1
		fmt.Printf("[%v] %s\n", u, v.Name())
	}

	for {
		v, err := components.ReadInt("Enter: ")
		if err != nil {
			return nil, fmt.Errorf("failed userinput: %w", err)
		}
		v--

		if v < 0 || v >= len(sites.Sites) {
			continue
		}

		return sites.Sites[v], nil
	}
}

func Search(site model.Site) (result model.SearchResult, err error) {
	for {
		search, err := components.ReadString("Search: ")
		if err != nil {
			return model.SearchResult{}, fmt.Errorf("failed userinput: %w", err)
		}

		results, err := site.Search(search)
		if err != nil {
			log.Error("Failed to search", "service", site.Name(), "term", search)
			continue
		}

		if len(results) == 0 {
			continue
		}

		if len(results) == 1 {
			return results[0], nil
		}

		var values []huh.Option[int]
		for i, v := range results {
			values = append(values, huh.NewOption(v.Name, i))
		}

		var v int
		err = huh.NewSelect[int]().
			Title("Pick a season").
			Options(values...).
			Value(&v).
			WithTheme(components.Theme).
			Run()
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Selected season", "seasonLabel", results[v].Name)
		return results[v], nil
	}
}
