package components

import (
	"fmt"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

	"charm.land/log/v2"
)

func Language(site model.Site, episode model.Episode) (lang string, err error) {
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
		log.Error("no language found", "site", site, "href", episode.Href)
		return "", fmt.Errorf("no languages available")
	}

	if len(languages) == 1 {
		log.Info("Selected the only language available", "lang", languages[0])
		return languages[0], nil
	}

	for {
		for i := range languages {
			fmt.Printf("[%v] %s\n", i+1, languages[i])
		}

		u, err := ReadInt(Reader, "Select: ")
		if err != nil {
			return "", fmt.Errorf("failed userinput: %s", err)
		}
		u--

		if u < 0 || u >= len(languages) {
			util.ClearTerminal()
			log.Error("Invalid selection")
			continue
		}

		return languages[u], nil
	}
}
