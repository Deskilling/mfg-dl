package aniworld

import (
	"mfg-dl/internal/sites/modules/aniworld"
	"path/filepath"
	"testing"

	"github.com/Deskilling/gopkg/pkg/filesystem"
)

const episodeFiles = "./data/episodes"

func TestAniworldStreamParse(t *testing.T) {
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

		streams, err := aniworld.ParseStreams(string(content))
		if err != nil {
			t.Errorf("failed parsing episodes %s: %v", v.Name(), err)
			continue
		}

		for i, u := range streams {
			if u.Href == "" {
				t.Errorf("%s [%d]: empty href", v.Name(), i)
			}
			if u.Hoster == "" {
				t.Errorf("%s [%d]: empty hoster", v.Name(), i)
			}
			if u.Language == "" {
				t.Errorf("%s [%d]: empty language", v.Name(), i)
			}
		}

		t.Logf("Passed for %s (%d Streams)", v.Name(), len(streams))
	}
}
