package endpoint

import (
	"encoding/json"
	"net/http"

	"mfg-dl/internal/sites/model"

	"charm.land/log/v2"
)

/*
 * POST /download
 * expects json of model.Stream
 * returns status
 */
func HandleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var stream model.Stream
	err := json.NewDecoder(r.Body).Decode(&stream)
	if err != nil {
		writeError(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	site := getSite(stream.Service)
	if site == nil {
		writeError(w, "service not found", http.StatusBadRequest)
	}

	err = site.Download(stream)
	if err != nil {
		log.Error("failed episodes", "err", err)
		writeError(w, "episodes failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
