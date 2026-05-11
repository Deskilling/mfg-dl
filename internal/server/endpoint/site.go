package endpoint

import (
	"encoding/json"
	"net/http"

	"charm.land/log/v2"
)

/*
 * Get /searchsite?q=...&service=...
 * q = query
 * service = service
 * returns []model.SearchResults as json
 */
func HandleSearchSite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, "missing query parameter: q", http.StatusBadRequest)
		return
	}

	service := r.URL.Query().Get("service")
	if service == "" {
		writeError(w, "missing query parameter: service", http.StatusBadRequest)
		return
	}

	site := getSite(service)
	if site == nil {
		writeError(w, "service not found", http.StatusBadRequest)
	}

	result, err := site.Search(query)
	if err != nil {
		log.Errorf("failed search", "err", err)
		writeError(w, "search failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}
