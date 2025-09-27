package aniworld

type Endpoints map[string]string

// to avoid isp blocking this could also be the updated ip
// TODO - curl --http2 --header "accept: application/dns-json" "https://one.one.one.one/dns-query?name=example.com"
var BaseURL = "https://aniworld.to"

var AniEndpoints = Endpoints{
	"default":  BaseURL,
	"search":   BaseURL + "/ajax/seriesSearch?keyword=",
	"episodes": BaseURL + "/anime/stream/",
}

type Languages map[string]string

var AniLanguages = Languages{
	"1": "gerdub",
	"2": "engdub",
	"3": "gersub",
}
