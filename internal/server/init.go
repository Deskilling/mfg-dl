package server

import (
	"log/slog"
	"net/http"
	"strconv"

	"mfg-dl/internal/core"
	"mfg-dl/internal/server/endpoint"
)

func Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", endpoint.HandleHealth)
	mux.HandleFunc("/services", endpoint.HandleServices)
	mux.HandleFunc("/search", endpoint.HandleSearch)
	mux.HandleFunc("/cover", endpoint.HandleCover)
	mux.HandleFunc("/score", endpoint.HandleScore)
	mux.HandleFunc("/searchsite", endpoint.HandleSearchSite)
	mux.HandleFunc("/seasons", endpoint.HandleSeasons)
	mux.HandleFunc("/episodes", endpoint.HandleEpisodes)
	mux.HandleFunc("/streams", endpoint.HandleStreams)
	mux.HandleFunc("/download", endpoint.HandleDownload)

	slog.Info("Starting server", "addr", core.GetConfig().Server.Port)
	return http.ListenAndServe(":"+strconv.Itoa(core.GetConfig().Server.Port), middleware(mux))
}
