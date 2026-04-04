package m3u

// based on https://github.com/jamesnetherton/m3u/

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type VariantStream struct {
	Resolution       string
	Bandwidth        int
	AverageBandwidth int
	Codecs           string
	Name             string
	FrameRate        float64
	HdcpLevel        string
	Video            string
	Audio            string
	Subtitle         string
	ClosedCaptions   string
	URI              string
}

// TODO Improve error handle stuff and read how m3u actually works cheffe €€€€ :(
func Parse(f io.ReadCloser) (variantStream []VariantStream, err error) {
	defer f.Close()

	firstLine := true
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if firstLine && !strings.HasPrefix(line, "#EXTM3U") {
			return nil, errors.New(("invalid m3u file format. Expected #EXTM3U file header"))
		}

		firstLine = false

		if strings.HasPrefix(line, "#EXTINF") {
			line := strings.ReplaceAll(line, "#EXTINF:", "")
			trackInfo := strings.Split(line, ",")
			if len(trackInfo) < 2 {
				return nil, errors.New("invalid m3u file format. Expected EXTINF metadata to contain track length and name data")
			}
		} else if strings.HasPrefix(line, "#EXT-X-STREAM-INF") {
			line := strings.ReplaceAll(line, "#EXT-X-STREAM-INF:", "")
			streamInfo := strings.Split(line, ",")
			if len(streamInfo) < 1 {
				return nil, errors.New(("invalid m3u file format. Expected EXT-X-STREAM-INF metadata to contain bitrate data"))
			}
			stream := &VariantStream{}
			for i, param := range streamInfo {
				if strings.HasPrefix(param, "BANDWIDTH") {
					bandwidth := strings.Split(streamInfo[i], "=")[1]
					bandwidthInt, err := strconv.Atoi(bandwidth)
					if err != nil {
						return nil, fmt.Errorf("unable to parse bandwidth: %w", err)
					}
					stream.Bandwidth = bandwidthInt
				}
				if strings.HasPrefix(param, "AVERAGE-BANDWIDTH") {
					averageBandwidth := strings.Split(streamInfo[i], "=")[1]
					averageBandwidthInt, err := strconv.Atoi(averageBandwidth)
					if err != nil {
						return nil, fmt.Errorf("unable to parse average bandwidth: %s", err)
					}
					stream.AverageBandwidth = averageBandwidthInt
				}
				if strings.HasPrefix(param, "CODECS") {
					codecs := strings.Split(streamInfo[i], "=")[1]
					stream.Codecs = codecs
				}
				if strings.HasPrefix(param, "RESOLUTION") {
					resolution := strings.Split(streamInfo[i], "=")[1]
					stream.Resolution = resolution
				}
				if strings.HasPrefix(param, "FRAME-RATE") {
					frameRate := strings.Split(streamInfo[i], "=")[1]
					frameRateFloat, err := strconv.ParseFloat(frameRate, 64)
					if err != nil {
						return nil, fmt.Errorf("unable to parse frame rate", err)
					}
					stream.FrameRate = frameRateFloat
				}
				if strings.HasPrefix(param, "HDCP-LEVEL") {
					hdcpLevel := strings.Split(streamInfo[i], "=")[1]
					stream.HdcpLevel = hdcpLevel
				}
				if strings.HasPrefix(param, "VIDEO") {
					video := strings.Split(streamInfo[i], "=")[1]
					stream.Video = video
				}
				if strings.HasPrefix(param, "AUDIO") {
					audio := strings.Split(streamInfo[i], "=")[1]
					stream.Audio = audio
				}
				if strings.HasPrefix(param, "SUBTITLES") {
					subtitle := strings.Split(streamInfo[i], "=")[1]
					stream.Subtitle = subtitle
				}
				if strings.HasPrefix(param, "CLOSED-CAPTIONS") {
					closedCaptions := strings.Split(streamInfo[i], "=")[1]
					stream.ClosedCaptions = closedCaptions
				}
				if strings.HasPrefix(param, "NAME") {
					name := strings.Split(streamInfo[i], "=")[1]
					stream.Name = name
				}
			}
			variantStream = append(variantStream, *stream)
		} else if strings.HasPrefix(line, "#") || line == "" {
			continue
		} else if variantStream != nil {
			variantStream[len(variantStream)-1].URI = strings.Trim(line, " ")
		}
	}

	return variantStream, nil
}
