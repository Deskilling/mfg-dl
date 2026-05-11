package server

import (
	"net/http"

	"mfg-dl/internal/server/endpoint"

	"charm.land/log/v2"
)

func Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", endpoint.HandleHealth)
	mux.HandleFunc("/services", endpoint.HandleServices)
	mux.HandleFunc("/search", endpoint.HandleSearch)
	mux.HandleFunc("/score", endpoint.HandleScore)
	mux.HandleFunc("/searchsite", endpoint.HandleSearchSite)
	mux.HandleFunc("/seasons", endpoint.HandleSeasons)
	mux.HandleFunc("/episodes", endpoint.HandleEpisodes)
	mux.HandleFunc("/streams", endpoint.HandleStreams)
	mux.HandleFunc("/download", endpoint.HandleDownload)

	log.Info("Starting server", "addr", ":6702")
	return http.ListenAndServe(":6702", middleware(mux))
}
