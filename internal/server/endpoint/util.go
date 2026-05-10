package endpoint

import (
	"encoding/json"
	"net/http"

	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
)

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(apiError{Code: code, Message: message})
}

func getSite(name string) model.Site {
	for _, v := range sites.Sites {
		if v.Name() == name {
			return v
		}
	}

	return nil
}
