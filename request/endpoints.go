package request

type Endpoints map[string]string

var BaseURL = "https://aniworld.to"

var AniworldEndpoints = Endpoints{
	"default":  BaseURL,
	"search":   BaseURL + "/ajax/seriesSearch?keyword=",
	"episodes": BaseURL + "/anime/stream/",
}
