package server

import (
	"encoding/json"
	"net/http"

	"mfg-dl/internal/search"
	"mfg-dl/internal/sites"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/util"

	"charm.land/log/v2"
)

func Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/search", handleSearch)
	mux.HandleFunc("/match", handleMatch)
	mux.HandleFunc("/searchSite", handleSearchSite)
	mux.HandleFunc("/seasons", handleSeasons)
	mux.HandleFunc("/episodes", handleEpisodes)
	mux.HandleFunc("/streams", handleStreams)
	mux.HandleFunc("/download", handleDownloads)

	log.Info("Starting server", "addr", ":6702")
	return http.ListenAndServe(":6702", middleware(mux))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}

	results, err := search.Search(query)
	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}

func handleMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing name", http.StatusBadRequest)
		return
	}

	selected := model.SearchResult{
		Service: "TMDB",
		Name:    name,
	}

	var allServiceResults [][]model.SearchResult
	for _, site := range sites.Sites {
		results, err := site.Search(util.NormalizeString(name))
		if err != nil {
			log.Error("search failed", "service", site.Name(), "err", err)
			continue
		}
		if len(results) > 0 {
			allServiceResults = append(allServiceResults, results)
		}
	}

	search.Match(&selected, allServiceResults)

	json.NewEncoder(w).Encode(selected)
}

func handleSearchSite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}

	serviceName := r.URL.Query().Get("service")
	if serviceName == "" {
		http.Error(w, "missing service", http.StatusBadRequest)
		return
	}

	var site model.Site
	for _, s := range sites.Sites {
		if s.Name() == serviceName {
			site = s
			break
		}
	}

	if site == nil {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}

	result, _ := site.Search(query)

	json.NewEncoder(w).Encode(result)
}

func handleSeasons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	serviceName := r.URL.Query().Get("service")
	if serviceName == "" {
		http.Error(w, "missing service", http.StatusBadRequest)
		return
	}

	var season model.Season
	err := json.NewDecoder(r.Body).Decode(&season)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	var site model.Site
	for _, s := range sites.Sites {
		if s.Name() == serviceName {
			site = s
			break
		}
	}

	if site == nil {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}

	episodes, err := site.Episodes(season)
	if err != nil {
		log.Error("failed getting episodes", "service", serviceName, "err", err)
		http.Error(w, "failed getting episodes", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(episodes)
}

func handleEpisodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func handleStreams(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleDownloads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
