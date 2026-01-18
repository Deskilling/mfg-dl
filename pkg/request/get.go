package request

import (
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
)

func Get(endpoint string) (string, error) {
	log.Debug("Sending get request to", "url", endpoint)
	req, err := NewRequest(http.MethodGet, endpoint)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		log.Error(err)
		return "", err
	}

	resp, err := Client.Do(req)
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
