package serienstream

import (
	"testing"

	"mfg-dl/internal/sites/modules/serienstream"
)

var terms []string = []string{"breaking bad", "fallout", "squid game"}

func TestSerienstreamSearchRequest(t *testing.T) {
	site := serienstream.New()
	for _, v := range terms {
		results, err := site.Search(v)
		if err != nil {
			t.Fatalf("failed search for %q: %v", v, err)
		}

		t.Logf("Passed for %s (%d results)", v, len(results))
	}
}
