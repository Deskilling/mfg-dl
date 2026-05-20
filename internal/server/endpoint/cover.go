package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"mfg-dl/internal/core"
	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/model"

	"github.com/Deskilling/gopkg/pkg/filesystem"
)

/*
 * POST /cover
 * expects json of model.searchResult
 * returns path to cover
 */
func HandleCover(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var searchResult model.SearchResult
	err := json.NewDecoder(r.Body).Decode(&searchResult)
	if err != nil {
		writeError(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("%s/cache/images/%s-%s.tmp", filepath.Join(*filesystem.GetExecPath(), core.GetConfig().Location.Temp), searchResult.Service, searchResult.Name)
	if !filesystem.ExistPath(path) {
		request.DownloadFile(nil, searchResult.Cover, path)
	}

	w.Write([]byte(path))
}
