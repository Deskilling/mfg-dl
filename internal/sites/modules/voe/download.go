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

func (service Voe) BaseDownload(voeUrl, output string) (err error) {
	if filesystem.ExistPath(output) {
		return nil
	}

	baseHtml, err := request.Get(service.client, voeUrl)
	if err != nil {
		return fmt.Errorf("Failed to get voe: %w", err)
	}

	baseUrl, err := VoeUrlHtml(string(baseHtml))
	if err != nil {
		return fmt.Errorf("Failed to extract voe url: %w", err)
	}

	err = service.PlayerDownload(baseUrl, output)
	if err != nil {
		return fmt.Errorf("failed player download: %w", err)
	}

	return nil
}

func (service Voe) PlayerDownload(voeUrl, output string) (err error) {
	if filesystem.ExistPath(output) {
		return nil
	}

	voeHtml, err := request.Get(service.client, voeUrl)
	if err != nil {
		return fmt.Errorf("failed getting voe player: %w", err)
	}

	parsed, err := Parse(string(voeHtml))
	if err != nil {
		return fmt.Errorf("failed parsing voe streams: %w", err)
	}

	parsed.Directory = util.RemoveAfterSymbol(parsed.FileCode, "/")

	if core.GetConfig().Downloads.FfmpegDownload {
		err = ffmpeg.DownloadHLS(parsed.Source, output)
		if err != nil {
			return fmt.Errorf("failed downloading hls stream with ffmpeg: %w", err)
		}
		return nil
	}

	masterTxt, err := request.Get(service.client, parsed.Source)
	if err != nil {
		return fmt.Errorf("failed get master file: %w", err)
	}

	master, err := m3u.Parse(io.NopCloser(strings.NewReader(string(masterTxt))))
	if err != nil {
		return fmt.Errorf("failed parsing master file: %w", err)
	}

	baseUrl, err := GetBaseUrl(parsed.Source)
	if err != nil {
		return fmt.Errorf("failed getting voe baseurl: %w", err)
	}

	indexTxt, err := request.Get(service.client, baseUrl+master[0].URI)
	if err != nil {
		return fmt.Errorf("failed get index file: %w", err)
	}

	index, err := m3u.ParseIndex(io.NopCloser(strings.NewReader(string(indexTxt))))
	if err != nil {
		return fmt.Errorf("failed parsing index: %w", err)
	}

	dir := fmt.Sprintf("%s/segments/%s/", core.GetConfig().Location.Temp, parsed.Directory)
	err = m3u.DownloadSegments(index, baseUrl, dir)
	if err != nil {
		return fmt.Errorf("failed to download all segments for %s: %w", dir, err)
	}

	err = ffmpeg.ConvertTSFilesToVideo(dir, output)
	if err != nil {
		// TODO THIS ALSO MEANS IT FAILED IMPLEMENTATION BEFORE WAS KINDA ASS
		return fmt.Errorf("failed creating video out of segments: %w", err)

	}

	return nil
}

// Uses best available quality
// TODO
func GetBaseUrl(input string) (url string, err error) {
	re := regexp.MustCompile(`(.*?)/[^/]*\.m3u8`)
	match := re.FindStringSubmatch(input)

	if len(match) == 0 {
		return "", fmt.Errorf("no base url found in: %s", input)
	}

	return match[1] + "/", nil
}
