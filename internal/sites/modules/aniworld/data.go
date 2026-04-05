package aniworld

import (
	"net/http"
	"time"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/modules/voe"
)

const Name = "Aniworld"
const BaseURL = "https://aniworld.to"

type Endpoints map[string]string

var AniEndpoints = Endpoints{
	"default":  BaseURL,
	"search":   BaseURL + "/ajax/seriesSearch?keyword=",
	"episodes": BaseURL + "/anime/stream/",
}

type Aniworld struct {
	client *http.Client
	voe    *voe.Voe
}

var Hoster []string = []string{"VOE"}

type Languages map[string]string

var AniLanguages = Languages{
	"1": "gerdub",
	"2": "engsub",
	"3": "gersub",
}

func New() *Aniworld {
	return &Aniworld{
		client: &http.Client{
			Timeout:   0 * time.Second,
			Transport: request.Client.Transport,
		},
		voe: voe.New(),
	}
}

func (site *Aniworld) Name() (service string) {
	return Name
}
