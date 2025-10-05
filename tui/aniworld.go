package tui

import (
	"fmt"
	"strconv"
	"strings"

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
	if err != nil {
		log.Error(err)
		return
	}

	if len(seasons) == 0 {
		log.Error("no seasons found")
		return
	}

	var filme bool = false
	fmt.Println("Available Seasons:")
	for i, v := range seasons {
		if seasons[0].Label == "Alle Filme" {
			filme = true
			fmt.Printf("[%v] %s\n", i, v.Label)
		} else {
			fmt.Printf("[%v] %s\n", i+1, v.Label)
		}
	}

	input = GetUserInput("Enter seasons (e.g., 1 2 3 or all): ")
	var selectedSeasons []string

	if strings.ToLower(strings.TrimSpace(input)) == "all" {
		for i := range seasons {
			selectedSeasons = append(selectedSeasons, strconv.Itoa(i))
		}
	} else {
		parts := strings.Fields(input)
		for _, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				log.Errorf("invalid input: %s is not a number", part)
				return
			}

			var seasonIndex int
			if filme {
				if num < 0 || num >= len(seasons) {
					log.Errorf("invalid selection: %d", num)
					return
				}
				seasonIndex = num
			} else {
				if num < 1 || num > len(seasons) {
					log.Errorf("invalid selection: %d", num)
					return
				}
				seasonIndex = num - 1
			}
			selectedSeasons = append(selectedSeasons, strconv.Itoa(seasonIndex))
		}
	}

	if len(selectedSeasons) == 0 {
		log.Error("no seasons selected")
		return
	}

	firstSelectedSeason := selectedSeasons[0]
	streams, err := aniworld.GetStreams(anime, firstSelectedSeason, "1")
	if err != nil {
		log.Warn("Could not pre-determine languages, you may encounter issues.", "err", err)
	}

	var languages []string
	var language string
	for _, v := range streams {
		if !util.Contains(languages, v.Language) {
			languages = append(languages, v.Language)
		}
	}

	if len(languages) >= 2 {
		for i, v := range languages {
			fmt.Printf("[%v] %s\n", i+1, aniworld.AniLanguages[v])
		}

		input = GetUserInput("Select Language: ")
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
	} else if len(languages) == 1 {
		language = languages[0]
	} else {
		log.Fatal("No Languages found")
	}

	log.Infof("Starting download for %d season(s)", len(selectedSeasons))
	aniworld.DownloadSeason(anime, language, "VOE", selectedSeasons)
}
