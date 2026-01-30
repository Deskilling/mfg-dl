package request

import (
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	clients = make(map[string]*http.Client)
	mu      sync.Mutex
)

func Get(url string, proxyAddr string) (content string, err error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	if !(proxyAddr == "") {
		client, err = getSock5Client(proxyAddr)
		if err != nil {
			return "", err
		}
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
