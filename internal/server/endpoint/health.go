package endpoint

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

/*
 * GET OR POST /health
 * returns status
 */
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	slog.Error("failed encoding json", "err", err)

}
