package serienstream

import (
	"path/filepath"
	"strings"
	"testing"

	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/sites/modules/serienstream"

	"github.com/Deskilling/gopkg/pkg/filesystem"
)

const serienFiles = "./data/stream"

func TestSerienstreamSeasonParse(t *testing.T) {
	files, err := filesystem.ReadDirectory(serienFiles, ".txt")
	if err != nil {
		t.Fatalf("failed reading %s directory: %v", serienFiles, err)
	}

	for _, v := range files {
		content, err := filesystem.ReadFile(filepath.Join(serienFiles, v.Name()))
		if err != nil {
			t.Errorf("failed reading file %s: %v", v.Name(), err)
			continue
		}

		seasons, err := serienstream.ParseSeasons(string(content))
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
		Name: "Breaking Bad",
		Href: "/serie/breaking-bad",
	},
	{
		Name: "Squid Game",
		Href: "/serie/squid-game",
	},
	{
		Name: "Fallout",
		Href: "/serie/fallout",
	},
}

func TestSerienstreamSeasonLive(t *testing.T) {
	site := serienstream.New()

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
