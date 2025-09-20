package parser

// based on https://github.com/jamesnetherton/m3u/

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

type VariantStream struct {
	Resolution      string
	Bandwidth       int
	AverageBandwith int
	Codecs          string
	Name            string
	FrameRate       float64
	HdcpLevel       string
	Video           string
	Audio           string
	Subtitle        string
	ClosedCaptions  string
	URI             string
}

func Parse(fileName string) ([]VariantStream, error) {
	var f io.ReadCloser

	if strings.HasPrefix(fileName, "http://") || strings.HasPrefix(fileName, "https://") {
		data, err := http.Get(fileName)
		if err != nil {
			err = fmt.Errorf("unable to open playlist URL: %v", err)
			log.Error(err)
			return nil, err
		}
		f = data.Body
	} else {
		file, err := os.Open(fileName)
		if err != nil {
			err = fmt.Errorf("unable to open playlist file: %v", err)
			log.Error(err)
			return nil, err
		}
		f = file
	}
	defer f.Close()

	var variantStream []VariantStream

	firstLine := true
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if firstLine && !strings.HasPrefix(line, "#EXTM3U") {
			err := errors.New("invalid m3u file format. Expected #EXTM3U file header")
			log.Error(err)
			return nil, err
		}

		firstLine = false

		if strings.HasPrefix(line, "#EXTINF") {
			line := strings.Replace(line, "#EXTINF:", "", -1)
			trackInfo := strings.Split(line, ",")
			if len(trackInfo) < 2 {
				err := errors.New("invalid m3u file format. Expected EXTINF metadata to contain track length and name data")
				log.Error(err)
				return nil, err
			}
		} else if strings.HasPrefix(line, "#EXT-X-STREAM-INF") {
			line := strings.Replace(line, "#EXT-X-STREAM-INF:", "", -1)
			streamInfo := strings.Split(line, ",")
			if len(streamInfo) < 1 {
				err := errors.New("invalid m3u file format. Expected EXT-X-STREAM-INF metadata to contain bitrate data")
				log.Error(err)
				return nil, err
			}
			stream := &VariantStream{}
			for i, param := range streamInfo {
				if strings.HasPrefix(param, "BANDWIDTH") {
					bandwidth := strings.Split(streamInfo[i], "=")[1]
					bandwidthInt, err := strconv.Atoi(bandwidth)
					if err != nil {
						err = fmt.Errorf("unable to parse bandwidth: %w", err)
						log.Error(err)
						return nil, err
					}
					stream.Bandwidth = bandwidthInt
				}
				if strings.HasPrefix(param, "AVERAGE-BANDWIDTH") {
					averageBandwidth := strings.Split(streamInfo[i], "=")[1]
					averageBandwidthInt, err := strconv.Atoi(averageBandwidth)
					if err != nil {
						err = fmt.Errorf("unable to parse average bandwidth: %w", err)
						log.Error(err)
						return nil, err
					}
					stream.AverageBandwith = averageBandwidthInt
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
						err = fmt.Errorf("unable to parse frame rate: %w", err)
						log.Error(err)
						return nil, err
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
