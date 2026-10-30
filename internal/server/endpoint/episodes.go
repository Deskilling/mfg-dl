package endpoint

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"mfg-dl/internal/sites/model"
)

/*
 * POST /episodes
 * expects json of model.Season
 * returns []model.Episode as json
 */
func HandleEpisodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var season model.Season
	err := json.NewDecoder(r.Body).Decode(&season)
	if err != nil {
		writeError(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	site := getSite(season.Service)
	if site == nil {
		writeError(w, "service not found", http.StatusBadRequest)
	}

	episodes, err := site.Episodes(season)
	if err != nil {
		slog.Error("failed episodes", "err", err)
		writeError(w, "episodes failed", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(episodes)
	slog.Error("failed encoding json", "err", err)

}
