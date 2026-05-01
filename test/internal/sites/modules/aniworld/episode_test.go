package aniworld

import (
	"path/filepath"
	"strings"
	"testing"

	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/sites/modules/aniworld"

	"github.com/Deskilling/gopkg/pkg/filesystem"
)

const seasonFiles = "./data/stream"

func TestAniworldEpisodeParse(t *testing.T) {
	files, err := filesystem.ReadDirectory(seasonFiles, ".txt")
	if err != nil {
		t.Fatalf("failed reading %s directory: %v", seasonFiles, err)
	}

	for _, v := range files {
		content, err := filesystem.ReadFile(filepath.Join(seasonFiles, v.Name()))
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

// SeasonNum && SeasonLabel dont matter in this context, only the Href is relevant
var sampleSeasons []model.Season = []model.Season{
	{
		Name:        "Detektiv Conan",
		Href:        "/anime/stream/detektiv-conan/staffel-1",
		SeasonNum:   "01",
		SeasonLabel: "Staffel 1",
	},
	{
		Name:        "My Neighbor Totoro",
		Href:        "/anime/stream/my-neighbor-totoro/staffel-1",
		SeasonNum:   "01",
		SeasonLabel: "Staffel 1",
	},
	{
		Name:        "One Piece",
		Href:        "/anime/stream/one-piece/staffel-1",
		SeasonNum:   "01",
		SeasonLabel: "Staffel 1",
	},
	{
		Name:        "Princess Mononoke",
		Href:        "/anime/stream/princess-mononoke/staffel-1",
		SeasonNum:   "01",
		SeasonLabel: "Staffel 1",
	},
	{
		Name:        "Rascal Does Not Dream of Bunny Girl Senpai",
		Href:        "/anime/stream/rascal-does-not-dream-of-bunny-girl-senpai/staffel-1",
		SeasonNum:   "01",
		SeasonLabel: "Staffel 1",
	},
}

func TestAniworldEpisodeLive(t *testing.T) {
	site := aniworld.New()

	for _, v := range sampleSeasons {
		seasons, err := site.Episodes(v)
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
				t.Errorf("%s [%d]: empty href", s.Name, i)
			}
			if strings.Contains(s.EpisodeNum, " ") {
				t.Errorf("%s [%d]: invalid episode number: %q", s.Name, i, s.EpisodeNum)
			}
			if s.EpisodeTitle == "" {
				t.Errorf("%s [%d]: invalid episode title: %q", s.Name, i, s.EpisodeTitle)
			}
		}
		t.Logf("Passed for %s (%d seasons)", v.Name, len(seasons))
	}
}
