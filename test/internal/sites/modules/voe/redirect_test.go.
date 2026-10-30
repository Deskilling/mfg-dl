package voe

import (
	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/modules/voe"
	"path/filepath"
	"testing"

	"github.com/Deskilling/gopkg/pkg/filesystem"
)

const redirectFiles = "./data/redirect"

func TestVoeRedirect(t *testing.T) {
	files, err := filesystem.ReadDirectory(redirectFiles, ".txt")
	if err != nil {
		t.Fatalf("failed reading %s directory: %v", redirectFiles, err)
	}

	for _, v := range files {
		content, err := filesystem.ReadFile(filepath.Join(redirectFiles, v.Name()))
		if err != nil {
			t.Errorf("failed reading file %s: %v", v.Name(), err)
			continue
		}

		url, err := voe.VoeRedirect(string(content))
		if err != nil {
			t.Errorf("failed regex on %s: %v", v.Name(), err)
			continue
		}

		t.Logf("Passed for %s: %s", v.Name(), url)
	}
}

var sampleVoeRedirectUrls []string = []string{
	"https://voe.sx/e/o6pq69uxu7wl",
	"https://voe.sx/e/6q1ngfkyvebf",
	"https://voe.sx/e/xqo6271z5tev",

	"https://voe.sx/e/lutc9vijlalf",
	"https://voe.sx/e/wark5ec0rq4w",
	"https://voe.sx/e/kukevkqehqyg",

	"https://voe.sx/e/dfwfjaxudbea",
}

func TestVoeRedirectLive(t *testing.T) {
	for _, v := range sampleVoeRedirectUrls {
		baseHtml, err := request.Get(nil, v)
		if err != nil {
			t.Errorf("failed to get voe: %v", err)
			continue
		}

		_, err = voe.VoeRedirect(string(baseHtml))
		if err != nil {
			t.Errorf("failed regex on %s: %v", v, err)
			continue
		}

		t.Logf("Passed for %s", v)
	}
}
