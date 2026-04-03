package m3u

// based on https://github.com/jamesnetherton/m3u/

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
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

// TODO Improve error handle stuff and read how m3u actually works cheffe €€€€ :(
func ParseIndex(f io.ReadCloser) (m3u8Index Index, err error) {
	defer f.Close()

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
				if err == nil {
					m3u8Index.Version = version
				}
			} else if strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE") {
				sequenceStr := strings.TrimPrefix(line, "#EXT-X-MEDIA-SEQUENCE:")
				sequence, err := strconv.Atoi(sequenceStr)
				if err == nil {
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
			}
		}
	}

	err = scanner.Err()
	if err != nil {
		return
	}

	if m3u8Index.Segments == nil {
		err = errors.New("no segments found in index")
	}

	return
}
