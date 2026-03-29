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
)

func BaseDownload(voeUrl, output string) (err error) {
	if filesystem.ExistPath(output) {
		return nil
	}

	baseHtml, err := request.Get(voeUrl)
	if err != nil {
		return err
	}

	baseUrl, err := VoeUrlHtml(string(baseHtml))
	if err != nil {
		return err
	}

	PlayerDownload(baseUrl, output)

	return nil
}

func PlayerDownload(voeUrl, output string) (err error) {
	if filesystem.ExistPath(output) {
		return nil
	}

	voeHtml, err := request.Get(voeUrl)
	if err != nil {
		return err
	}

	parsed, err := Parse(string(voeHtml))
	if err != nil {
		return err
	}

	// make sure its chill
	parsed.Directory = util.RemoveAfterSymbol(parsed.FileCode, "/")

	if core.GetConfig().Downloads.FfmpegDownload {
		ffmpeg.DownloadHLS(parsed.Source, output)
		return nil
	}

	masterTxt, err := request.Get(parsed.Source)
	if err != nil {
		return err
	}

	master, err := m3u.Parse(io.NopCloser(strings.NewReader(string(masterTxt))))
	if err != nil {
		return err
	}

	baseUrl := GetBaseUrl(parsed.Source)

	indexTxt, err := request.Get(baseUrl + master[0].URI)
	if err != nil {
		return err
	}

	index, err := m3u.ParseIndex(io.NopCloser(strings.NewReader(string(indexTxt))))
	if err != nil {
		return err
	}

	dir := fmt.Sprintf("%s/segments/%s/", core.GetConfig().Location.Temp, parsed.Directory)
	err = m3u.DownloadSegments(index, baseUrl, dir)
	if err != nil {
		return fmt.Errorf("failed to download all segments for %s", dir)
	}

	err = ffmpeg.ConvertTSFilesToVideo(dir, output)
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
