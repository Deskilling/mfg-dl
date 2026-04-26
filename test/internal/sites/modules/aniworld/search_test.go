package aniworld

import (
	"testing"

	"mfg-dl/internal/sites/modules/aniworld"
)

var terms []string = []string{"a", "totoro", "rascal", "conan", "one piece"}

func TestAniworldSearchRequest(t *testing.T) {
	site := aniworld.New()
	for _, v := range terms {
		results, err := site.Search(v)
		if err != nil {
			t.Fatalf("failed search for %q: %v", v, err)
		}

		t.Logf("Passed for %s (%d results)", v, len(results))
	}
}
