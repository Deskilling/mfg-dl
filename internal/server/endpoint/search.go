package endpoint

import (
	"encoding/json"
	"net/http"

	"mfg-dl/internal/search"

	"charm.land/log/v2"
)

/*
 * GET /search?q=...
 * q = query
 * Searches TMDB
 * returns json object of []model.SearchResult
 */
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, "missing query parameter: q", http.StatusBadRequest)
		return
	}

	results, err := search.Search(query)
	if err != nil {
		log.Error("search failed", "query", query, "err", err)
		writeError(w, "search failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}
