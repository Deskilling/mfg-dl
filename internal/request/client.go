package request

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"mfg-dl/internal/core"
	"mfg-dl/internal/util"
	"mfg-dl/pkg/filesystem"

	"charm.land/log/v2"
)

var userAgent string = "deskilling/mfg-dl"

var Client = &http.Client{
	Timeout: 0,
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,

		ForceAttemptHTTP2:     true,
		MaxIdleConns:          512,
		MaxIdleConnsPerHost:   256,
		MaxConnsPerHost:       0,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func cachePath(endpoint string) string {
	base := core.GetConfig().Location.Temp + "/cache/requests/"
	if core.GetConfig().Debug.Sha256Cache {
		return base + util.Hash256String(endpoint)
	}
	return base + endpoint
}

func Get(endpoint string, headers ...map[string]string) ([]byte, error) {
	if core.GetConfig().Cache.EnableCache {
		path := cachePath(endpoint)

		if filesystem.ExistPath(path) {
			data, err := filesystem.ReadFile(path)
			if err == nil {
				log.Debugf("cache hit for %s", endpoint)
				return []byte(data), nil
			}
		}
	}

	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("User-Agent", userAgent)

	for _, h := range headers {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}

	resp, err := Client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("status %d", resp.StatusCode, "url", endpoint)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	log.Debugf("request to %s took: %v", endpoint, time.Since(start))

	if core.GetConfig().Cache.EnableCache {
		path := cachePath(endpoint)

		filesystem.WriteFile(path, body)
	}

	return body, nil
}
