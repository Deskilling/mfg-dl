package components

import (
	"fmt"
	"mfg-dl/internal/sites/model"
)

func Language(site *model.Site, episode model.Episode) (lang string, err error) {
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
		u := ReadInt(Reader, "Select: ")
		u--

		if u < 0 || u >= len(languages) {
			continue
		}

		return languages[u], nil
	}
}
