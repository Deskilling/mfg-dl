package voe

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"mfg-dl/internal/core"
	"mfg-dl/internal/ffmpeg"
	"mfg-dl/internal/m3u"
	"mfg-dl/internal/request"
	"mfg-dl/internal/util"
	"mfg-dl/pkg/filesystem"

	"github.com/charmbracelet/log"
)

func BaseDownload(voeUrl, output string) (err error) {
	if filesystem.ExistPath(output) {
		log.Info("Already downloaded", "output", output)
		return nil
	}

	baseHtml, err := request.Get(voeUrl)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debug("got html for", "voeUrl", voeUrl)

	baseUrl, err := VoeUrlHtml(baseHtml)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debug("got baseurl from html", "baseUrl", baseUrl)

	PlayerDownload(baseUrl, output)

	return nil
}

func PlayerDownload(voeUrl, output string) (err error) {
	if filesystem.ExistPath(output) {
		log.Info("Already downloaded", "output", output)
		return nil
	}

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

	if core.GetConfig().Extra.FfmpegDownload {
		ffmpeg.DownloadHLS(parsed.Source, output)
		return nil
	}

	masterTxt, err := request.Get(parsed.Source)
	if err != nil {
		log.Error(err)
		return err
	}

	master, err := m3u.Parse(io.NopCloser(strings.NewReader(masterTxt)))
	if err != nil {
		log.Error(err)
		return err
	}

	baseUrl := GetBaseUrl(parsed.Source)
	log.Debug("baseurl", "baseUrl", baseUrl+master[0].URI)

	indexTxt, err := request.Get(baseUrl + master[0].URI)
	if err != nil {
		log.Error(err)
		return err
	}

	index, err := m3u.ParseIndex(io.NopCloser(strings.NewReader(indexTxt)))
	if err != nil {
		log.Error(err)
		return err
	}

	dir := fmt.Sprintf("%s/%s/", core.GetConfig().Location.Temp, parsed.Directory)
	log.Info(dir)
	err = m3u.DownloadSegments(index, baseUrl, dir)
	if err != nil {
		return fmt.Errorf("failed to download all segments for %s", dir)
	}

	err = m3u.ConvertTSFilesToVideo(dir, output)
	if err != nil {
		// TODO THIS ALSO MEANS IT FAILED IMPLEMENTATION BEFORE WAS KINDA ASS
		return err

	}

	return nil
}

// Uses best available quality
// TODO
func GetBaseUrl(input string) string {
	re := regexp.MustCompile(`(.*?)/[^/]*\.m3u8`)
	match := re.FindStringSubmatch(input)

	if len(match) <= 0 {
		return ""
	}
	return match[1] + "/"
}
