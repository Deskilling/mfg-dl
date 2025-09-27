package voe

import (
	"regexp"

	"mfg-dl/filesystem"
	"mfg-dl/m3u"
	"mfg-dl/request"
	"mfg-dl/util"

	"github.com/charmbracelet/log"
)

func BaseDownload(voeUrl, output string) error {
	baseHtml, err := request.Get(voeUrl)
	if err != nil {
		log.Error(err)
		return err
	}

	baseUrl, err := VoeUrlHtml(baseHtml)
	if err != nil {
		log.Error(err)
		return err
	}

	err = PlayerDownload(baseUrl, output)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func PlayerDownload(voeUrl, output string) error {
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
		m3u.ConvertTSFilesToVideo(filesystem.GetExecDir()+"/temp/"+parsed.Directory, output)
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
