package voe

import (
	"mfg-dl/m3u"
	"mfg-dl/request"
	"mfg-dl/util"
	"regexp"

	"github.com/charmbracelet/log"
)

func Download(voeUrl, output string) error {
	voeHtml, err := request.Get(voeUrl)
	if err != nil {
		log.Error(err)
		return err
	}

	parsed, err := Parse(voeHtml)
	if err != nil {
		log.Error(err)
		return err
	}

	// make sure its chill
	parsed.Directory = util.RemoveAfterSymbol(parsed.FileCode, "/")

	master, err := m3u.Parse(parsed.Source)
	if err != nil {
		log.Error(err)
		return err
	}

	baseUrl := GetBaseUrl(parsed.Source)
	log.Debug("baseurl", "baseUrl", baseUrl+master[0].URI)

	index, err := m3u.ParseIndex(baseUrl + master[0].URI)
	if err != nil {
		log.Error(err)
		return err
	}

	if m3u.DownloadSegments(index, baseUrl, "./temp/"+parsed.Directory+"/") {
		log.Debug("yes")
		m3u.ConvertTSFilesToVideo("./temp/"+parsed.Directory, output)
	}

	return nil
}

func GetBaseUrl(input string) string {
	re := regexp.MustCompile(`(.*?)/[^/]*\.m3u8`)
	match := re.FindStringSubmatch(input)

	if len(match) <= 0 {
		return ""
	}
	return match[1] + "/"
}
