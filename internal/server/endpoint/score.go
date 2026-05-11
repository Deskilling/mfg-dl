package endpoint

import (
	"encoding/json"
	"net/http"

	"mfg-dl/internal/search"
	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

	"charm.land/log/v2"
)

/*
 * POST /score
 * expects json of model.SearchResult
 * returns model.SearchResult as json
 */
func HandleScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var selected model.SearchResult
	if err := json.NewDecoder(r.Body).Decode(&selected); err != nil {
		writeError(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if selected.Name == "" {
		writeError(w, "missing field: name", http.StatusBadRequest)
		return
	}

	var allServiceResults [][]model.SearchResult
	for _, site := range sites.Sites {
		results, err := site.Search(util.NormalizeString(selected.Name))
		if err != nil {
			log.Error("score search failed", "service", site.Name(), "err", err)
			continue
		}
		if len(results) > 0 {
			allServiceResults = append(allServiceResults, results)
		}
	}

	if len(allServiceResults) == 0 {
		writeError(w, "no results found on any service", http.StatusNotFound)
		return
	}

	search.Match(&selected, allServiceResults)
	json.NewEncoder(w).Encode(selected)
}
