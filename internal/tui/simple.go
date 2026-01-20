package tui

import (
	"fmt"
	"mfg-dl/internal/sites"

	"github.com/charmbracelet/log"
)

func SimpleTui() {
	var userNum int = 0
	if len(sites.Sites) != 1 {
		for i, v := range sites.Sites {
			u := i + 1
			fmt.Printf("[%v] %s\n", u, v.Name)
		}

		fmt.Print("Enter: ")
		_, err := fmt.Scanf("%d", &userNum)
		if err != nil {
			log.Error("failed reading user input", "err", err)
			return
		}
		userNum--
	}

	var userSearch string
	fmt.Print("Search: ")
	_, err := fmt.Scanf("%s", &userSearch)
	if err != nil {
		log.Error("failed reading user input", "err", err)
		return
	}

	search, err := sites.Sites[userNum].Search(userSearch)
	if err != nil {
		log.Error("failed getting sites", "input", userSearch, "err", err)
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

	seasons, err := sites.Sites[userNum].Seasons(search[userSelected])
	if err != nil {
		log.Error("failed getting seasons", "input", search[userSelected], "err", err)
	}
	for i, v := range seasons {
		if seasons[0].Label == "Alle Filme" {
			fmt.Printf("[%v] %s\n", i, v.Label)
		} else {
			u := i + 1
			fmt.Printf("[%v] %s\n", u, v.Label)
		}
	}

	var userSeason int
	fmt.Print("Search: ")
	_, err = fmt.Scanf("%v", &userSeason)
	if err != nil {
		log.Error("failed reading user input", "err", err)
		return
	}

	if seasons[0].Label != "Alle Filme" {
		userSeason--
	}

	episode, err := sites.Sites[userNum].Episodes(seasons[userSeason])
	if err != nil {
		log.Error("failed getting seasons", "input", search[userSelected], "err", err)
	}
	for i, v := range episode {
		u := i + 1
		fmt.Printf("[%v] %s\n", u, v.Title)
	}

	var userEpisode int
	fmt.Print("Search: ")
	_, err = fmt.Scanf("%v", &userEpisode)
	if err != nil {
		log.Error("failed reading user input", "err", err)
		return
	}
	userEpisode--

	stream, _ := sites.Sites[userNum].Streams(episode[userEpisode])
	if err != nil {
		log.Error("failed getting streams", "input", episode[userEpisode], "err", err)
	}

	for _, v := range stream {
		fmt.Println(v.Language)
	}
}
