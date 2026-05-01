package aniworld

import (
	"path/filepath"
	"strings"
	"testing"

	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/sites/modules/aniworld"

	"github.com/Deskilling/gopkg/pkg/filesystem"
)

const animeFiles = "./data/stream"

func TestAniworldSeasonParse(t *testing.T) {
	files, err := filesystem.ReadDirectory(animeFiles, ".txt")
	if err != nil {
		t.Fatalf("failed reading %s directory: %v", animeFiles, err)
	}

	for _, v := range files {
		content, err := filesystem.ReadFile(filepath.Join(animeFiles, v.Name()))
		if err != nil {
			t.Errorf("failed reading file %s: %v", v.Name(), err)
			continue
		}

		seasons, err := aniworld.ParseSeasons(string(content))
		if err != nil {
			t.Errorf("failed parsing seasons %s: %v", v.Name(), err)
			continue
		}

		for i, u := range seasons {
			if u.Href == "" {
				t.Errorf("%s [%d]: empty href", v.Name(), i)
			}
			if u.SeasonLabel == "" {
				t.Errorf("%s [%d]: empty season label", v.Name(), i)
			}
			if strings.Contains(u.SeasonNum, " ") {
				t.Errorf("%s [%d]: invalid season number: %q", v.Name(), i, u.SeasonNum)
			}
		}

		t.Logf("Passed for %s (%d Seasons)", v.Name(), len(seasons))
	}
}

var sampleSearchResults []model.SearchResult = []model.SearchResult{
	{
		Name: "Detektiv Conan",
		Href: "/anime/stream/detektiv-conan",
	},
	{
		Name: "My Neighbor Totoro",
		Href: "/anime/stream/my-neighbor-totoro",
	},
	{
		Name: "One Piece",
		Href: "/anime/stream/one-piece",
	},
	{
		Name: "Princess Mononoke",
		Href: "/anime/stream/princess-mononoke",
	},
	{
		Name: "Rascal Does Not Dream of Bunny Girl Senpai",
		Href: "/anime/stream/rascal-does-not-dream-of-bunny-girl-senpai",
	},
}

func TestAniworldSeasonLive(t *testing.T) {
	site := aniworld.New()

	for _, v := range sampleSearchResults {
		seasons, err := site.Seasons(v)
		if err != nil {
			t.Errorf("failed getting seasons for %s: %v", v.Name, err)
			continue
		}
		if len(seasons) == 0 {
			t.Errorf("%s: expected seasons, got none", v.Name)
			continue
		}

		for i, s := range seasons {
			if s.Href == "" {
				t.Errorf("%s [%d]: empty href", v.Name, i)
			}
			if s.SeasonLabel == "" {
				t.Errorf("%s [%d]: empty season label", v.Name, i)
			}
			if strings.Contains(s.SeasonNum, " ") {
				t.Errorf("%s [%d]: invalid season number: %q", v.Name, i, s.SeasonNum)
			}
		}
		t.Logf("Passed for %s (%d seasons)", v.Name, len(seasons))
	}
}
