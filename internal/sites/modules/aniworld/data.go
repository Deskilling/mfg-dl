package aniworld

import "mfg-dl/internal/sites/model"

var BaseURL = "https://aniworld.to"

type Endpoints map[string]string

var AniEndpoints = Endpoints{
	"default":  BaseURL,
	"search":   BaseURL + "/ajax/seriesSearch?keyword=",
	"episodes": BaseURL + "/anime/stream/",
}

var Site = model.Site{
	Name:    "Aniworld",
	Baseurl: BaseURL,

	Search:           GetSearch,
	Seasons:          GetSeasons,
	Episodes:         GetEpisodes,
	Streams:          GetStreams,
	DownloadMultiple: DownloadMultiple,
}

var Hoster []string = []string{"VOE"}

type Languages map[string]string

var AniLanguages = Languages{
	"1": "gerdub",
	"2": "engsub",
	"3": "gersub",
}
