package aniworld

import (
	"mfg-dl/internal/sites/model"
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

// SeasonNum && EpisodeTitle... dont matter in this context, only the Href is relevant
var sampleEpisodes []model.Episode = []model.Episode{
	{
		Name:       "Detektiv Conan",
		Href:       "/anime/stream/detektiv-conan/staffel-1/episode-1",
		SeasonNum:  "01",
		EpisodeNum: "01",
	},
	{
		Name:       "My Neighbor Totoro",
		Href:       "/anime/stream/my-neighbor-totoro/staffel-1/episode-1",
		SeasonNum:  "01",
		EpisodeNum: "01",
	},
	{
		Name:       "One Piece",
		Href:       "/anime/stream/one-piece/staffel-1/episode-1",
		SeasonNum:  "01",
		EpisodeNum: "01",
	},
	{
		Name:       "Princess Mononoke",
		Href:       "/anime/stream/princess-mononoke/staffel-1/episode-1",
		SeasonNum:  "01",
		EpisodeNum: "01",
	},
	{
		Name:       "Rascal Does Not Dream of Bunny Girl Senpai",
		Href:       "/anime/stream/rascal-does-not-dream-of-bunny-girl-senpai/staffel-1/episode-1",
		SeasonNum:  "01",
		EpisodeNum: "01",
	},
	{
		Name:       "Rascal Does Not Dream of Bunny Girl Senpai",
		Href:       "/anime/stream/rascal-does-not-dream-of-bunny-girl-senpai/filme/film-1",
		SeasonNum:  "00",
		EpisodeNum: "00",
	},
}

func TestAniworldStreamLive(t *testing.T) {
	site := aniworld.New()

	for _, v := range sampleEpisodes {
		seasons, err := site.Streams(v)
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
			if s.Hoster == "" {
				t.Errorf("%s [%d]: empty hoster", s.Name, i)
			}
			if s.Language == "" {
				t.Errorf("%s [%d]: empty language", s.Name, i)
			}
		}
		t.Logf("Passed for %s %s %s (%d streams)", v.Name, v.SeasonNum, v.EpisodeNum, len(seasons))
	}
}
