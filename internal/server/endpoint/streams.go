package endpoint

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"mfg-dl/internal/sites/model"
)

/*
 * POST /streams
 * expects json of model.Episode
 * returns []model.Stream as json
 */
func HandleStreams(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var episode model.Episode
	err := json.NewDecoder(r.Body).Decode(&episode)
	if err != nil {
		writeError(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	site := getSite(episode.Service)
	if site == nil {
		writeError(w, "service not found", http.StatusBadRequest)
	}

	streams, err := site.Streams(episode)
	if err != nil {
		slog.Error("failed streams", "err", err)
		writeError(w, "streams failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(streams)
}
