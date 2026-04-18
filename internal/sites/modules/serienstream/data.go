package serienstream

import (
	"net/http"
	"time"

	"mfg-dl/internal/request"
	"mfg-dl/internal/sites/modules/voe"
)

const Name = "Serienstream"
const BaseURL = "http://186.2.175.5"

type Endpoints map[string]string

var SerienstreamEndpoints = Endpoints{
	"default": BaseURL,
	"search":  BaseURL + "/api/search/suggest?term=",
}

type Serienstream struct {
	client *http.Client
	voe    *voe.Voe
}

var Hoster []string = []string{"VOE"}

type Languages map[string]string

var SerienstreamLanguages = Languages{
	"1": "gerdub",
	"2": "engsub",
	"3": "gersub",
}

func New() *Serienstream {
	return &Serienstream{
		client: &http.Client{
			Timeout:   0 * time.Second,
			Transport: request.Client.Transport,
		},
		voe: voe.New(),
	}
}

func (site *Serienstream) Name() (service string) {
	return Name
}
