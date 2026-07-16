package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"mfg-dl/internal/core"
	"mfg-dl/internal/request"
	"mfg-dl/internal/search"
	"mfg-dl/internal/sites/model"

	"charm.land/log/v2"
	"github.com/Deskilling/gopkg/pkg/filesystem"
)

/*
 * GET OR POST /search?q=...
 * q = query
 * Searches TMDB
 * returns json object of []model.SearchResult
 */
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, "missing query parameter: q", http.StatusBadRequest)
		return
	}

	results, err := search.Search(query)
	if err != nil {
		log.Error("search failed", "query", query, "err", err)
		writeError(w, "search failed", http.StatusInternalServerError)
		return
	}

	cache := r.URL.Query().Get("cacheCover")
	if cache != "" {
		for i := range results {
			results[i].CoverPath = filepath.Join(
				*filesystem.GetExecPath(),
				core.GetConfig().Location.Temp,
				"cache",
				"images",
				fmt.Sprintf("%s-%s.tmp", results[i].Service, results[i].Name),
			)
		}

		go func(results []model.SearchResult) {
			for _, v := range results {
				if !filesystem.ExistPath(v.CoverPath) {
					err := request.DownloadFile(nil, v.Cover, v.CoverPath)
					if err != nil {
						log.Error("failed to cache image", "url", v.Cover, "path", v.CoverPath, "err", err)
					}
				}
			}
		}(results)
	}
	json.NewEncoder(w).Encode(results)
}
