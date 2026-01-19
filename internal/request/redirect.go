package request

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
)

func Redirect(link string) (string, error) {
	client := &http.Client{
		Timeout:       Client.Timeout,
		Transport:     Client.Transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	req, err := NewRequest(http.MethodGet, link)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		log.Error(err)
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil && resp == nil {
		err = fmt.Errorf("failed request of redirect link: %w", err)
		log.Error(err)
		return "", err
	}
	if resp == nil {
		err = fmt.Errorf("no response for redirect link")
		log.Error(err)
		return "", err
	}
	defer resp.Body.Close()

	location, err := resp.Location()
	if err != nil {
		err = fmt.Errorf("no location for response %w", err)
		log.Error(err)
		return "", err
	}

	log.Debug("Found redirect url", "source", link, "to", location.String())
	return location.String(), nil
}
