package request

import (
	"net/http"
)

func Redirect(link string) (string, error) {
	client := Client

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	req, err := http.NewRequest("HEAD", link, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil && resp == nil {
		return "", err
	}
	defer resp.Body.Close()

	httpLocation, err := resp.Location()
	if err != nil {
		return "", err
	}

	return httpLocation.String(), nil
}
