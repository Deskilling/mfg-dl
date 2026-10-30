package endpoint

import (
	"encoding/json"
	"log/slog"
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
	err := json.NewEncoder(w).Encode(apiError{Code: code, Message: message})
	slog.Error("failed encoding json", "err", err)

}

func getSite(name string) model.Module {
	for _, v := range sites.Modules {
		if v.Name() == name {
			return v
		}
	}

	return nil
}
