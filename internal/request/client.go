package request

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:          500,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func Get(endpoint string) (string, error) {
	start := time.Now()

	resp, err := client.Get(endpoint)
	if err != nil {
		err = fmt.Errorf("request failed: %w", err)
		log.Error(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("request failed with status code: %d", resp.StatusCode)
		log.Error(err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("reading body failed: %w", err)
		log.Error(err)
		return "", err
	}

	log.Debugf("request to %s took: %v", endpoint, time.Since(start))
	return string(body), nil
}
