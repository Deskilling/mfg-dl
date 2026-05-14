package endpoint

import (
	"encoding/json"
	"net/http"

	"mfg-dl/internal/sites"
)

/*
 * GET OR POST /services
 * returns list of enabled services
 */
func HandleServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var names []string
	for _, site := range sites.Sites {
		names = append(names, site.Name())
	}

	json.NewEncoder(w).Encode(names)
}
