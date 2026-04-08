package m3u

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"mfg-dl/pkg/filesystem"
	"mfg-dl/pkg/m3u8"
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

		streams, err := m3u8.ParseIndex(io.NopCloser(strings.NewReader(string(content))))
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
