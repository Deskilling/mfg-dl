package request

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

func Get(endpoint string) (string, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		log.Error(err)
		return "", err
	}

	req.Header.Set("User-Agent", "Deskilling/mfg-dl")

	resp, err := client.Do(req)
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

	return string(body), nil
}
