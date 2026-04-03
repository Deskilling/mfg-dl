package components

import (
	"fmt"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

	"charm.land/log/v2"
)

func Language(site *model.Site, episode model.Episode) (lang string, err error) {
	stream, err := site.Streams(episode)
	if err != nil {
		return
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
		err = fmt.Errorf("no languages available")
		return
	}

	for {
		for i := range languages {
			fmt.Printf("[%v] %s\n", i+1, languages[i])
		}

		var u int
		u, err = ReadInt(Reader, "Select: ")
		if err != nil {
			err = fmt.Errorf("failed userinput: %s", err)
			return
		}
		u--

		if u < 0 || u >= len(languages) {
			util.ClearTerminal()
			log.Error("Invalid selection")
			continue
		}

		lang = languages[u]
		return
	}
}
