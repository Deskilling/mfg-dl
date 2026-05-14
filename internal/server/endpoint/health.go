package endpoint

import (
	"encoding/json"
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

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
