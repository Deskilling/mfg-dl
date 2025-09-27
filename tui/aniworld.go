package tui

import (
	"fmt"
	"strconv"

	"mfg-dl/sites/aniworld"
	"mfg-dl/util"

	"github.com/charmbracelet/log"
)

func Aniworld() {
	results, err := aniworld.GetSearch(GetUserInput("Enter Anime: "))
	if err != nil {
		log.Error(err)
		return
	}

	for i, v := range results {
		fmt.Printf("[%v] %s\n", i+1, v.Name)
	}

	input := GetUserInput("Enter: ")
	userAnime, err := strconv.Atoi(input)
	if err != nil {
		log.Error(err)
		return
	}
	if userAnime < 1 || userAnime > len(results) {
		log.Error("invalid anime selection")
		return
	}

	anime := results[userAnime-1].Link

	seasons, err := aniworld.GetSeasons(anime)
	userAnime, _ = strconv.Atoi(input)

	anime = results[userAnime-1].Link

	seasons, err = aniworld.GetSeasons(anime)
	if err != nil {
		log.Error(err)
		return
	}

	// checks every time but who cares
	if len(seasons) == 0 {
		log.Error("no seasons found")
		return
	}

	for i, v := range seasons {
		if seasons[0].Label == "Alle Filme" {
			fmt.Printf("[%v] %s\n", i, v.Label)
		} else {
			fmt.Printf("[%v] %s\n", i+1, v.Label)
		}
	}

	input = GetUserInput("Enter: ")
	season := input

	episodes, err := aniworld.GetEpisodes(anime, input)
	if err != nil {
		log.Error(err)
		return
	}

	for i, v := range episodes {
		fmt.Printf("[%v] %s\n", i+1, v.Title)
	}

	input = GetUserInput("Enter: ")
	episode := input

	streams, err := aniworld.GetStreams(anime, season, episode)
	if err != nil {
		log.Error(err)
		return
	}
	var languages []string
	var language string

	for _, v := range streams {
		if !util.Contains(languages, v.Language) {
			languages = append(languages, v.Language)
			if len(languages) == 0 {
				log.Error("no streams available for the selected episode")
				return
			}

			var userLanguage int
			if len(languages) >= 2 {
				for i, v := range languages {
					log.Debug(v)
					fmt.Printf("[%v] %s\n", i+1, aniworld.AniLanguages[v])
				}

				input = GetUserInput("Enter: ")
				userLanguage, err := strconv.Atoi(input)
				if err != nil {
					log.Error(err)
					return
				}
				if userLanguage < 1 || userLanguage > len(languages) {
					log.Error("invalid language selection")
					return
				}
				language = languages[userLanguage-1]
			} else {
				language = languages[0]
			}
			language = languages[userLanguage-1]
		} else {
			language = languages[0]
		}

		aniworld.Download(anime, season, episode, language, "VOE")

	}
}
