package request

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
)

func Redirect(link string) (location string, err error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	req, err := http.NewRequest("HEAD", link, nil)
	if err != nil {
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

	httpLocation, err := resp.Location()
	if err != nil {
		err = fmt.Errorf("no location for response %w", err)
		log.Error(err)
		return "", err
	}

	location = httpLocation.String()

	log.Debug("Found redirect url", "source", link, "to", location)
	return location, nil
}
