package endpoint

import (
	"encoding/json"
	"net/http"
)

/*
 * GET /health
 * returns status
 */
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
