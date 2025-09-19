package request

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
)

func Redirect(link string) (string, error) {
	client := &http.Client{
		Timeout:       client.Timeout,
		Transport:     client.Transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		log.Error(err)
		return "", err
	}
	req.Header.Set("User-Agent", "Deskilling/mfg-dl")

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

	return location.String(), nil
}
