package request

import (
	"net/http"
)

func Redirect(link string) (url string, err error) {
	client := Client

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
	req, err := http.NewRequest("HEAD", link, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return
	}
	defer resp.Body.Close()

	httpLocation, err := resp.Location()
	if err != nil {
		return
	}

	url = httpLocation.String()
	return
}
