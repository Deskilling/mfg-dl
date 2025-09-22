package m3u

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

type Index struct {
	TargetDuration string
	AllowCache     bool
	PlaylistType   string
	Version        int
	Sequence       int
	Segments       []Segment
}

type Segment struct {
	Duration float64
	URI      string
}

func ParseIndex(filename string) (*Index, error) {
	var f io.ReadCloser

	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
		resp, err := http.Get(filename)
		if err != nil {
			err = fmt.Errorf("unable to open playlist URL: %w", err)
			log.Error(err)
			return nil, err
		}
		f = resp.Body
	} else {
		file, err := os.Open(filename)
		if err != nil {
			err = fmt.Errorf("unable to open playlist file: %w", err)
			log.Error(err)
			return nil, err
		}
		f = file
	}
	defer f.Close()

	m3u8Index := &Index{}
	var currentSegment *Segment

	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNum++

		if line == "" || strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "#EXTM3U") {
				continue
			}
			if strings.HasPrefix(line, "#EXT-X-TARGETDURATION") {
				durationStr := strings.TrimPrefix(line, "#EXT-X-TARGETDURATION:")
				m3u8Index.TargetDuration = durationStr
			} else if strings.HasPrefix(line, "#EXT-X-ALLOW-CACHE") {
				m3u8Index.AllowCache = strings.TrimPrefix(line, "#EXT-X-ALLOW-CACHE:") == "YES"
			} else if strings.HasPrefix(line, "#EXT-X-PLAYLIST-TYPE") {
				m3u8Index.PlaylistType = strings.TrimPrefix(line, "#EXT-X-PLAYLIST-TYPE:")
			} else if strings.HasPrefix(line, "#EXT-X-VERSION") {
				versionStr := strings.TrimPrefix(line, "#EXT-X-VERSION:")
				version, err := strconv.Atoi(versionStr)
				if err != nil {
					// kinda useless rn but log anyway
					log.Debug("could not parse version", "line", lineNum, "error", err)
				} else {
					m3u8Index.Version = version
				}
			} else if strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE") {
				sequenceStr := strings.TrimPrefix(line, "#EXT-X-MEDIA-SEQUENCE:")
				sequence, err := strconv.Atoi(sequenceStr)
				if err != nil {
					// all chill
					log.Debug("could not parse media sequence", "line", lineNum, "error", err)
				} else {
					m3u8Index.Sequence = sequence
				}
			} else if strings.HasPrefix(line, "#EXTINF") {
				info := strings.TrimPrefix(line, "#EXTINF:")
				parts := strings.SplitN(info, ",", 2)
				if len(parts) != 2 {
					continue
				}

				duration, err := strconv.ParseFloat(parts[0], 64)
				if err != nil {
					log.Warn("could not parse segment duration", "line", lineNum, "error", err)
					continue
				}

				currentSegment = &Segment{
					Duration: duration,
				}
			}
		} else {
			if currentSegment != nil {
				currentSegment.URI = line
				m3u8Index.Segments = append(m3u8Index.Segments, *currentSegment)
				currentSegment = nil
			} else {
				log.Warn("found a URI without a preceding EXTINF tag", "line", lineNum, "uri", line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if m3u8Index.Segments == nil {
		return nil, errors.New("no segments found in the playlist")
	}

	return m3u8Index, nil
}
