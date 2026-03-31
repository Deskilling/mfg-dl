package voe

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"mfg-dl/internal/util"

	"github.com/PuerkitoBio/goquery"
)

type VoeStream struct {
	Key                     string   `json:"key"`
	Sharing                 bool     `json:"sharing"`
	LogoEnabled             bool     `json:"logo_enabled"`
	LogoPath                string   `json:"logo_path"`
	LogoURL                 string   `json:"logo_url"`
	LogoPosition            string   `json:"logo_position"`
	Thumbnail               string   `json:"thumbnail"`
	ShowTitle               bool     `json:"show_title"`
	Airplay                 bool     `json:"airplay"`
	Check                   bool     `json:"check"`
	FileCode                string   `json:"file_code"`
	MetadataPreload         string   `json:"metadata_preload"`
	BufferLength            int      `json:"buffer_length"`
	BufferSize              int64    `json:"buffer_size"`
	DisableTimeSlider       bool     `json:"disable_timeslider"`
	Title                   string   `json:"title"`
	Source                  string   `json:"source"`
	Fallback                []string `json:"fallback"`
	Captions                []string `json:"captions"`
	DefaultCaptionsLanguage string   `json:"default_captions_language"`
	Request                 string   `json:"request"`
	DirectAccessAllowed     bool     `json:"direct_access_allowed"`
	DirectAccessURL         string   `json:"direct_access_url"`
	SDKVersion              string   `json:"sdk_version"`
	SiteName                string   `json:"site_name"`
	// Only used for the temp dir
	Directory string
}

func Parse(html string) (stream *VoeStream, err error) {
	if html == "" {
		return nil, fmt.Errorf("not html parsed")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("could not create goquery document: %w", err)
	}

	jsonElem := doc.Find("script[type='application/json']").First()
	if jsonElem.Length() == 0 {
		return nil, fmt.Errorf("no JSON found")
	}

	string := strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(jsonElem.Text()), `"]`), `["`)
	string = util.Rot13(string)
	string = VoeRemovePatterns(string)

	string, err = util.Base64Decode(string)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)

	}

	string = util.ShiftChars(string, 3)
	string = util.ReverseString(string)

	decoded, err := util.Base64Decode(string)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	replacer := strings.NewReplacer(`\/`, `/`)
	decoded = replacer.Replace(decoded)

	var data VoeStream
	err = json.Unmarshal([]byte(decoded), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to umasharl json: %w", err)
	}

	return &data, nil
}

func VoeRemovePatterns(str string) (result string) {
	patterns := []string{"@$", "^^", "~@", "%?", "*~", "!!", "#&"}
	for _, pat := range patterns {
		str = strings.ReplaceAll(str, pat, "")
	}
	return str
}

func VoeUrlHtml(htmlContent string) (url string, err error) {
	re := regexp.MustCompile(`window.location.href\s*=\s*['"](https://[^'"]+)['"]`)

	matches := re.FindStringSubmatch(htmlContent)

	if len(matches) <= 0 {
		return "", fmt.Errorf("no URL found in the provided HTML content")
	}

	return matches[1], nil
}
