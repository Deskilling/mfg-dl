package m3u

import (
	"io"
	"mfg-dl/internal/m3u"
	"mfg-dl/pkg/filesystem"
	"path/filepath"
	"strings"
	"testing"
)

const indexFiles string = "./data/index"

func TestIndexParser(t *testing.T) {
	files, err := filesystem.ReadDirectory(indexFiles, ".m3u")
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	for _, v := range files {
		content, err := filesystem.ReadFile(filepath.Join(indexFiles, v.Name()))
		if err != nil {
			t.Error(err)
			continue
		}

		streams, err := m3u.ParseIndex(io.NopCloser(strings.NewReader(string(content))))
		if err != nil {
			t.Errorf("failed to parse %s: %v", v.Name(), err)
			continue
		}

		if len(streams.Segments) == 0 {
			t.Errorf("expected segemnts in %s, got none", v.Name())
			continue
		}
	}
}
