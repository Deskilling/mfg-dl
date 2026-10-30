package voe

import (
	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/modules/voe"
	"path/filepath"
	"testing"

	"github.com/Deskilling/gopkg/pkg/filesystem"
)

const parseFiles = "./data/parse"

func TestVoeParse(t *testing.T) {
	files, err := filesystem.ReadDirectory(parseFiles, ".txt")
	if err != nil {
		t.Fatalf("failed reading %s directory: %v", parseFiles, err)
	}

	for _, v := range files {
		content, err := filesystem.ReadFile(filepath.Join(parseFiles, v.Name()))
		if err != nil {
			t.Errorf("failed reading file %s: %v", v.Name(), err)
			continue
		}

		parsed, err := voe.Parse(string(content))
		if err != nil {
			t.Errorf("failed parse on %s: %v", v.Name(), err)
			continue
		}

		if parsed.Source == "" {
			t.Errorf("failed parsing source on %s: %v", v.Name(), err)
			continue
		}

		baseUrl, err := voe.GetBaseUrl(parsed.Source)
		if err != nil {
			t.Errorf("failed parsing voe streams: %v", err)
			continue
		}

		if baseUrl == "" {
			t.Errorf("failed extracting url out of %s: %v", parsed.Source, err)
			continue
		}

		t.Logf("Passed for %s", v.Name())
	}
}

var sampleVoeParseUrl []string = []string{
	"https://jennifereconomicgive.com/e/o6pq69uxu7wl",
	"https://jennifereconomicgive.com/e/6q1ngfkyvebf",
	"https://jennifereconomicgive.com/e/xqo6271z5tev",
	"https://jennifereconomicgive.com/e/lutc9vijlalf",
	"https://jennifereconomicgive.com/e/wark5ec0rq4w",
	"https://jennifereconomicgive.com/e/kukevkqehqyg",
	"https://jennifereconomicgive.com/e/dfwfjaxudbea",
}

func TestVoeParseLive(t *testing.T) {
	for _, v := range sampleVoeParseUrl {
		voeHtml, err := request.Get(nil, v)
		if err != nil {
			t.Errorf("failed getting voe player: %v", err)
			continue
		}

		parsed, err := voe.Parse(string(voeHtml))
		if err != nil {
			t.Errorf("failed parsing voe streams: %v", err)
			continue
		}

		if parsed.Source == "" {
			t.Errorf("failed parsing source on %s: %v", v, err)
			continue
		}

		baseUrl, err := voe.GetBaseUrl(parsed.Source)
		if err != nil {
			t.Errorf("failed parsing voe streams: %v", err)
			continue
		}

		if baseUrl == "" {
			t.Errorf("failed extracting url out of %s: %v", parsed.Source, err)
			continue
		}

		t.Logf("Passed for %s", v)
	}
}
