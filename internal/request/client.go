package request

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"mfg-dl/internal/core"
	"mfg-dl/internal/util"
	"mfg-dl/pkg/filesystem"

	"charm.land/log/v2"
)

var userAgent string = "deskilling/mfg-dl"

var Client = &http.Client{
	Timeout: 30 * time.Second,
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

func Get(client *http.Client, endpoint string, headers ...map[string]string) (body []byte, err error) {
	if client == nil {
		client = Client
	}

	log.Debug("Sending Request", "url", endpoint)

	if core.GetConfig().Cache.EnableCache {
		path := cachePath(endpoint)

		if filesystem.ExistPath(path) {
			stat, err := os.Stat(path)
			if err == nil {
				if time.Since(stat.ModTime()) < time.Minute*time.Duration(core.GetConfig().Cache.Minutes) {
					log.Debug("Using cached request", "file", path)
					data, err := filesystem.ReadFile(path)
					if err == nil {
						return []byte(data), nil
					}
				}
			}
		}
	}

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
		return []byte{}, fmt.Errorf("status %d", resp.StatusCode)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if core.GetConfig().Cache.EnableCache {
		path := cachePath(endpoint)
		log.Debug("Writing to cache", "path", path)

		err = filesystem.WriteFile(path, body)
		if err != nil {
			log.Error("failed creating cache file", "path", path, "err", err)
		}
	}

	return body, nil
}
