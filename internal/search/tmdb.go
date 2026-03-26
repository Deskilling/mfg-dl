package search

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"
)

type SearchResult struct {
	Title  string
	Title2 string
	Title3 string
	Type   string
	Id     int
	Year   string
}

// var tmdbAPIKey = "TMDB_API_KEY"

type tmdbResponse struct {
	Results []tmdbResult `json:"results"`
}

type tmdbResult struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Name         string  `json:"name"`
	Overview     string  `json:"overview"`
	Popularity   float64 `json:"popularity"`
	ReleaseDate  string  `json:"release_date"`
	FirstAirDate string  `json:"first_air_date"`
	VoteAverage  float64 `json:"vote_average"`
	VoteCount    int     `json:"vote_count"`
	PosterPath   string  `json:"poster_path"`
	MediaType    string  `json:"media_type"`
}

func Search(query string) ([]model.SearchResult, error) {
	endpoint := fmt.Sprintf("https://www.themoviedb.org/search/multi?language=en&query=%s", url.QueryEscape(query))

	// https://www.themoviedb.org/search/multi?query=
	// https://api.themoviedb.org/3/search/multi?query=
	/*
		body, err := request.Get(endpoint, map[string]string{
			"Authorization": "Bearer " + tmdbAPIKey,
			"Accept":        "application/json",
		})
	*/

	body, err := request.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var result tmdbResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	sort.Slice(result.Results, func(i, j int) bool {
		return result.Results[i].Popularity > result.Results[j].Popularity
	})

	var results []model.SearchResult
	for _, v := range result.Results {
		if v.MediaType == "person" {
			continue
		}

		if v.VoteCount < 15 && v.Popularity < 1.0 {
			continue
		}

		if v.Title == "" {
			v.Title = v.Name
		}

		if v.ReleaseDate == "" {
			v.ReleaseDate = v.FirstAirDate
		}

		if len(v.ReleaseDate) >= 4 {
			v.ReleaseDate = v.ReleaseDate[:4]
		}

		results = append(results, model.SearchResult{
			Service:        "TMDB",
			Name:           util.NormalizeString(v.Title),
			Href:           fmt.Sprintf("/%s/%d", v.MediaType, v.ID),
			ProductionYear: v.ReleaseDate,
			Description:    v.Overview,
		})
	}

	return results, nil
}
