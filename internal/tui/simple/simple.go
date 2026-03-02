package simple

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

	"github.com/charmbracelet/log"
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

func SimpleTui() {
	site = SelectModule()
	result := Search()
	season := Seasons(result)
	streams := Episodes(season)

	if len(streams) == 0 {
		log.Error("no streams selected")
		return
	}

	site.DownloadMultiple(streams)
}

func SelectModule() *model.Site {
	if len(sites.Sites) == 1 {
		log.Infof("Selected %s", sites.Sites[0].Name)
		return &sites.Sites[0]
	}

	for i := range sites.Sites {
		u := i + 1
		fmt.Printf("[%v] %s\n", u, sites.Sites[i].Name)
	}

	for {
		v := readInt(reader, "Enter: ")
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
		search := readString(reader, "Search: ")

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
			log.Infof("Selected %s %s", results[0].Name, results[0].ProductionYear)
			return results[0]
		}

		for {
			for i := range results {
				u := i + 1
				fmt.Printf("[%v] %s %s\n", u, results[i].Name, results[i].ProductionYear)
			}

			v := readInt(reader, "Enter: ")
			v--

			if v < 0 || v >= len(results) {
				log.Error("Invalid Input, try again")
				continue
			}

			return results[v]
		}
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
