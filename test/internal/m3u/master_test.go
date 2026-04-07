package m3u

import (
	"io"
	"mfg-dl/internal/m3u"
	"mfg-dl/pkg/filesystem"
	"path/filepath"
	"strings"
	"testing"
)

const masterFiles string = "./data/master"

func TestMasterParser(t *testing.T) {
	files, err := filesystem.ReadDirectory(masterFiles, ".m3u8")
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	for _, v := range files {
		content, err := filesystem.ReadFile(filepath.Join(masterFiles, v.Name()))
		if err != nil {
			t.Error(err)
			continue
		}

		streams, err := m3u.Parse(io.NopCloser(strings.NewReader(string(content))))
		if err != nil {
			t.Errorf("failed to parse %s: %v", v.Name(), err)
			continue
		}

		if len(streams) == 0 {
			t.Errorf("expected variant streams in %s, got none", v.Name())
			continue
		}

		for _, stream := range streams {
			if stream.Bandwidth == 0 {
				t.Errorf("%s: stream has no bandwidth", v.Name())
			}

			if stream.URI == "" {
				t.Errorf("%s: stream has no URI", v.Name())
			}
		}
	}
}
