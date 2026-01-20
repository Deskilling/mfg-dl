package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
)

var cloudflaredomain string = "https://one.one.one.one/dns-query?name="

type Question struct {
	Name string
	Type int
}

type Answer struct {
	Name string
	Type int
	TTL  int
	Data string
}

type Cloudflareresponse struct {
	Status   int
	TC       bool
	RD       bool
	RA       bool
	AD       bool
	Question []Question
	Answer   []Answer
}

// TODO - curl --http2 --header "accept: application/dns-json" "https://one.one.one.one/dns-query?name=example.com"
func Cloudflaredns(domain string) (Cloudflareresponse, error) {
	client := &http.Client{
		Timeout: Client.Timeout,
	}

	req, err := http.NewRequest("GET", cloudflaredomain+domain, nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		log.Error(err)
		return Cloudflareresponse{}, err
	}
	req.Header.Add("User-Agent", "Deskilling/mfg-dl")
	req.Header.Add("accept", "application/dns-json")

	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("request failed: %w", err)
		log.Error(err)
		return Cloudflareresponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("request failed with status code: %d", resp.StatusCode)
		log.Error(err)
		return Cloudflareresponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("reading body failed: %w", err)
		log.Error(err)
		return Cloudflareresponse{}, err
	}

	var cloudflare Cloudflareresponse
	jsonbody := string(body)

	err = json.Unmarshal([]byte(jsonbody), &cloudflare)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal search results: %w", err)
		log.Error(err)
		return Cloudflareresponse{}, err
	}

	return cloudflare, nil
}
