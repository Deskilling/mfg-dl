package aniworld

type Endpoints map[string]string

var BaseURL = "https://aniworld.to"

var AniEndpoints = Endpoints{
	"default":  BaseURL,
	"search":   BaseURL + "/ajax/seriesSearch?keyword=",
	"episodes": BaseURL + "/anime/stream/",
}

type Languages map[string]string

var AniLanguages = Languages{
	"1": "german",
	"2": "engdub",
	"3": "gersub",
}
