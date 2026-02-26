package tui

import (
	"fmt"
	"mfg-dl/internal/sites"

	"github.com/charmbracelet/log"
)

func SimpleTui() {
	var userSite int = 0
	if len(sites.Sites) != 1 {
		for i, v := range sites.Sites {
			u := i + 1
			fmt.Printf("[%v] %s\n", u, v.Name)
		}

		fmt.Print("Enter: ")
		_, err := fmt.Scanf("%d", &userSite)
		if err != nil {
			log.Error("failed reading user input", "err", err)
			return
		}
		userSite--
	}

Search:
	var userSearch string
	fmt.Print("Search: ")
	_, err := fmt.Scanf("%s", &userSearch)
	if err != nil {
		log.Error("failed reading user input", "err", err)
		return
	}

	search, err := sites.Sites[userSite].Search(userSearch)
	if err != nil {
		log.Error("failed getting sites", "input", userSearch, "err", err)
		goto Search
	}
	for i, v := range search {
		u := i + 1
		fmt.Printf("[%v] %s\n", u, v.Name)
	}

	var userSelected int
	fmt.Print("Search: ")
	_, err = fmt.Scanf("%v", &userSelected)
	if err != nil {
		log.Error("failed reading user input", "err", err)
		return
	}
	userSelected--

	seasons, err := sites.Sites[userSite].Seasons(search[userSelected])
	if err != nil {
		log.Error("failed getting seasons", "input", search[userSelected], "err", err)
	}
	for i, v := range seasons {
		if seasons[0].SeasonLabel == "Alle Filme" {
			fmt.Printf("[%v] %s\n", i, v.SeasonLabel)
		} else {
			u := i + 1
			fmt.Printf("[%v] %s\n", u, v.SeasonLabel)
		}
	}

	var userSeason int
	fmt.Print("Search: ")
	_, err = fmt.Scanf("%v", &userSeason)
	if err != nil {
		log.Error("failed reading user input", "err", err)
		return
	}

	if seasons[0].SeasonLabel != "Alle Filme" {
		userSeason--
	}

	episode, err := sites.Sites[userSite].Episodes(seasons[userSeason])
	if err != nil {
		log.Error("failed getting seasons", "input", search[userSelected], "err", err)
	}
	for i, v := range episode {
		u := i + 1
		fmt.Printf("[%v] %s\n", u, v.EpisodeTitle)
	}

	var userEpisode int
	fmt.Print("Search: ")
	_, err = fmt.Scanf("%v", &userEpisode)
	if err != nil {
		log.Error("failed reading user input", "err", err)
		return
	}
	userEpisode--

	stream, _ := sites.Sites[userSite].Streams(episode[userEpisode])
	if err != nil {
		log.Error("failed getting streams", "input", episode[userEpisode], "err", err)
	}

	var languages []string
	var cnt uint = 0
	for _, v := range stream {
		seen := false
		for _, u := range languages {
			if u == v.Language {
				seen = true
				break
			}
		}

		if !seen {
			cnt++
			languages = append(languages, v.Language)
			fmt.Printf("[%v] %s\n", cnt, v.Language)
		}
	}

	var userLanguage int
	fmt.Print("Select: ")
	_, err = fmt.Scanf("%v", &userLanguage)
	if err != nil {
		log.Error("failed reading user input", "err", err)
		return
	}
	userLanguage--

	sites.Sites[userSite].Download(stream[0])
}
