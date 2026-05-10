package endpoint

import (
	"encoding/json"
	"net/http"

	"mfg-dl/internal/sites/model"

	"charm.land/log/v2"
)

func HandleEpisodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var searchResult model.Season
	err := json.NewDecoder(r.Body).Decode(&searchResult)
	if err != nil {
		writeError(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	site := getSite(searchResult.Service)
	if site == nil {
		writeError(w, "service not found", http.StatusBadRequest)
	}

	episodes, err := site.Episodes(searchResult)
	if err != nil {
		log.Errorf("failed episodes", "err", err)
		writeError(w, "episodes failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(episodes)
}
