package voe

import (
	"net/http"
	"time"

	"mfg-dl/internal/request"
)

type Voe struct {
	client *http.Client
}

func New() *Voe {
	return &Voe{
		client: &http.Client{
			Timeout:   0 * time.Second,
			Transport: request.Client.Transport,
		},
	}
}
