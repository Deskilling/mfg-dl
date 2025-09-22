package voe

import (
	"encoding/json"
	"fmt"
	"mfg-dl/util"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
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

func Parse(html string) (*VoeStream, error) {
	if html == "" {
		err := fmt.Errorf("not html parsed")
		log.Error(err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		err = fmt.Errorf("could not create goquery document: %w", err)
		log.Error(err)
		return nil, err
	}

	jsonElem := doc.Find("script[type='application/json']").First()
	if jsonElem.Length() == 0 {
		err = fmt.Errorf("no JSON found")
		return nil, err
	}

	string := strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(jsonElem.Text()), `"]`), `["`)
	string = util.Rot13(string)
	string = VoeRemovePatterns(string)

	string, err = util.Base64Decode(string)
	if err != nil {
		err = fmt.Errorf("failed to decode base64: %s", err)
		log.Error(err)
		return nil, err
	}

	string = util.ShiftChars(string, 3)
	string = util.ReverseString(string)

	decoded, err := util.Base64Decode(string)
	if err != nil {
		err = fmt.Errorf("failed to decode base64: %s", err)
		log.Error(err)
		return nil, err
	}

	replacer := strings.NewReplacer(`\/`, `/`)
	decoded = replacer.Replace(decoded)

	var data VoeStream
	err = json.Unmarshal([]byte(decoded), &data)
	if err != nil {
		err = fmt.Errorf("failed to umasharl json: %w", err)
		log.Error(err)
		return nil, err
	}

	return &data, nil
}

func VoeRemovePatterns(str string) string {
	patterns := []string{"@$", "^^", "~@", "%?", "*~", "!!", "#&"}
	result := str
	for _, pat := range patterns {
		result = strings.ReplaceAll(result, pat, "")
	}
	return result
}

// some people hate regex, but for this goquery is a bit overkill (and slower prob)
func VoeUrlHtml(htmlContent string) (string, error) {
	re := regexp.MustCompile(`window.location.href\s*=\s*['"](https://[^'"]+)['"]`)

	matches := re.FindStringSubmatch(htmlContent)

	if len(matches) <= 0 {
		err := fmt.Errorf("no URL found in the provided HTML content")
		log.Error(err)
		return "", err
	}

	return matches[1], nil
}
