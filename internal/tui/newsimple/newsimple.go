package newsimple

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"mfg-dl/internal/search"
	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

	"charm.land/log/v2"
)

var reader *bufio.Reader = bufio.NewReader(os.Stdin)
var site *model.Site

func readString(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Error("failed reading user input", "err", err)
		return ""
	}
	return strings.TrimSpace(input)
}

func readInt(reader *bufio.Reader, prompt string) int {
	for {
		input := readString(reader, prompt)
		val, err := strconv.Atoi(input)
		if err != nil {
			log.Error("invalid number input", "err", err)
			continue
		}
		return val
	}
}

func Tui() {
	tmdbResult := Search()
	service, index := Score(tmdbResult)

	site = &sites.Sites[index]
	result, _ := sites.Sites[index].Search(tmdbResult.Score.Query[service])
	season := Seasons(result[0])

	streams := Episodes(season)
	if len(streams) == 0 {
		log.Error("no streams selected")
		return
	}

	site.DownloadMultiple(streams)
}

func Search() model.SearchResult {
	for {
		input := readString(reader, "Search: ")

		tmdbResults, err := search.Search(input)
		if err != nil {
			log.Error("search failed", "err", err)
			continue
		}
		if len(tmdbResults) == 0 {
			log.Error("Not Found, try again")
			continue
		}

		selected := selectFromList(tmdbResults)

		var allServiceResults [][]model.SearchResult
		for _, site := range sites.Sites {
			results, err := site.Search(util.NormalizeString(selected.Name))
			if err != nil {
				log.Error(err)
				continue
			}
			if len(results) == 0 {
				results, err = site.Search(util.ShortSearchTerm(util.NormalizeString(selected.Name)))
				if len(results) == 0 {
					continue
				}
			}
			allServiceResults = append(allServiceResults, results)
		}

		search.Match(&selected, allServiceResults)
		return selected
	}
}

func selectFromList(results []model.SearchResult) model.SearchResult {
	if len(results) == 1 {
		log.Infof("Selected %s", results[0].Name)
		return results[0]
	}
	for {
		for i, v := range results {
			fmt.Printf("[%v] %s %s\n", i+1, v.Name, v.Href)
		}
		input := readInt(reader, "Enter: ")
		input--
		if input < 0 || input >= len(results) {
			log.Error("Invalid Input, try again")
			continue
		}
		return results[input]
	}
}

func Score(result model.SearchResult) (service string, index int) {
	if len(sites.Sites) == 1 {
		log.Error("fick dich hs")
		return sites.Sites[0].Service, 0
	}

	for {
		for i, v := range sites.Sites {
			fmt.Printf("[%v] %s %.2f%%\n", i+1, v.Service, result.Score.Score[v.Service]*100)
		}
		input := readInt(reader, "Enter: ")
		input--

		if input < 0 || input >= len(sites.Sites) {
			log.Error("Invalid Input, try again")
			continue
		}

		return sites.Sites[input].Service, input
	}
}

func Seasons(result model.SearchResult) model.Season {
	seasons, err := site.Seasons(result)
	if err != nil {
		log.Error("failed loading seasons", "err", err)
		return model.Season{}
	}

	for {
		if len(seasons) == 1 {
			log.Infof("Selected %s\n", seasons[0].SeasonLabel)
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

		v := readInt(reader, "Select: ")
		if seasons[0].SeasonLabel != "Alle Filme" {
			v--
		}

		if v < 0 || v >= len(seasons) {
			log.Error("Invalid Input, try again")
			continue
		}

		return seasons[v]
	}
}

func Episodes(season model.Season) []model.Stream {
	episodes, err := site.Episodes(season)
	if err != nil {
		log.Error("failed loading episodes", "err", err)
		return nil
	}

	for i := range episodes {
		if strings.TrimSpace(episodes[i].EpisodeTitle) == "" {
			episodes[i].EpisodeTitle = episodes[i].EpisodeAlternativeTitle
		}
	}

	var u []int

	if len(episodes) == 1 {
		log.Infof("Selected %s", episodes[0].EpisodeTitle)
		u = append(u, 0)
	} else {
		for {
			for i := range episodes {
				fmt.Printf("[%v] %s\n", i+1, episodes[i].EpisodeTitle)
			}

			v := readString(reader, "Select (use 1,2,3,4 or all for selection): ")

			if v == "all" {
				u = make([]int, len(episodes))
				for i := range episodes {
					u[i] = i
				}
				break
			}

			valid := []int{}
			for _, w := range util.ParseMultipleInts(v) {
				w--
				if w >= 0 && w < len(episodes) {
					valid = append(valid, w)
				}
			}

			if len(valid) > 0 {
				u = valid
				break
			}

			log.Error("invalid selection, please try again")
		}
	}

	lang, err := Language(episodes[u[0]])
	if err != nil {
		log.Error("language selection failed", "err", err)
		return nil
	}

	var unc []model.Stream
	for _, v := range u {
		stream, err := site.Streams(episodes[v])
		if err != nil {
			log.Error("failed loading streams", "err", err)
			continue
		}

		for _, w := range stream {
			if w.Language == lang {
				unc = append(unc, w)
				break
			}
		}
	}

	return unc
}

func Language(episode model.Episode) (lang string, err error) {
	stream, err := site.Streams(episode)
	if err != nil {
		return "", err
	}

	seen := make(map[string]struct{})
	var languages []string

	for _, v := range stream {
		if _, ok := seen[v.Language]; !ok {
			seen[v.Language] = struct{}{}
			languages = append(languages, v.Language)
		}
	}

	if len(languages) == 0 {
		return "", fmt.Errorf("no languages available")
	}

	for i := range languages {
		fmt.Printf("[%v] %s\n", i+1, languages[i])
	}

	for {
		u := readInt(reader, "Select: ")
		u--

		if u < 0 || u >= len(languages) {
			log.Error("Invalid Input, try again")
			continue
		}

		return languages[u], nil
	}
}
