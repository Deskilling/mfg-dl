package components

import (
	"fmt"
	"mfg-dl/internal/sites/model"

	"charm.land/huh/v2"
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

	var values []huh.Option[int]
	for i, v := range languages {
		values = append(values, huh.NewOption(v, i))
	}

	var v int
	err = huh.NewSelect[int]().
		Title("Pick a language").
		Options(values...).
		Value(&v).
		WithTheme(Theme).
		Run()
	if err != nil {
		return "", fmt.Errorf("failed userinput: %w", err)
	}

	return languages[v], nil
}
