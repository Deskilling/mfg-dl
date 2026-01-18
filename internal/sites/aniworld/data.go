package aniworld

type Endpoints map[string]string

var BaseURL = "https://aniworld.to"

var AniEndpoints = Endpoints{
	"default":  BaseURL,
	"search":   BaseURL + "/ajax/seriesSearch?keyword=",
	"episodes": BaseURL + "/anime/stream/",
}

var Hoster []string = []string{"VOE"}

type Languages map[string]string

var AniLanguages = Languages{
	"1": "gerdub",
	"2": "engsub",
	"3": "gersub",
}
