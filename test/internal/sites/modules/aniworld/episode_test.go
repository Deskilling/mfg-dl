package aniworld

import (
	"path/filepath"
	"strings"
	"testing"

	"mfg-dl/internal/sites/modules/aniworld"

	"github.com/Deskilling/gopkg/pkg/filesystem"
)

const episodeFiles = "./data/stream"

func TestAniworldEpisodeParse(t *testing.T) {
	files, err := filesystem.ReadDirectory(episodeFiles, ".txt")
	if err != nil {
		t.Fatalf("failed reading %s directory: %v", episodeFiles, err)
	}

	for _, v := range files {
		content, err := filesystem.ReadFile(filepath.Join(episodeFiles, v.Name()))
		if err != nil {
			t.Errorf("failed reading file %s: %v", v.Name(), err)
			continue
		}

		episodes, err := aniworld.ParseEpisodes(string(content))
		if err != nil {
			t.Errorf("failed parsing episodes %s: %v", v.Name(), err)
			continue
		}

		for i, u := range episodes {
			if u.Href == "" {
				t.Errorf("%s [%d]: empty href", v.Name(), i)
			}
			if strings.Contains(u.EpisodeNum, " ") {
				t.Errorf("%s [%d]: invalid episode number: %q", v.Name(), i, u.EpisodeNum)
			}
			if u.EpisodeTitle == "" {
				t.Errorf("%s [%d]: invalid episode title: %q", v.Name(), i, u.EpisodeTitle)
			}
		}

		t.Logf("Passed for %s (%d Episodes)", v.Name(), len(episodes))
	}
}
