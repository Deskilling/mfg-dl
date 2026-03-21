package request

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"mfg-dl/internal/core"
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

func Get(endpoint string) ([]byte, error) {
	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := Client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	log.Debugf("request to %s took: %v", endpoint, time.Since(start))

	if core.GetConfig().Debug.DumpHtml {
		filesystem.WriteFile(core.GetConfig().Location.Temp+"/requests/"+endpoint, body)
	}

	return body, nil
}
