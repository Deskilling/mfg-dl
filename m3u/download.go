package m3u

import (
	"mfg-dl/request"
	"strconv"
	"sync"

	"github.com/charmbracelet/log"
)

var maxConcurrency = 64

func DownloadSegments(index *Index, baseURL, directory string) bool {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrency)

	var failedDownloads []string
	var failedDirectories []string
	for i, v := range index.Segments {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(i int, v Segment) {
			defer wg.Done()
			defer func() { <-semaphore }()

			log.Debug("Downloading", "segment", v.URI)
			s := strconv.Itoa(i)
			err := request.DownloadFile(baseURL+v.URI, directory+s+".ts")
			if err != nil {
				failedDownloads = append(failedDownloads, baseURL+v.URI)
				failedDirectories = append(failedDirectories, directory+s+".ts")
				log.Error(err)
			}
		}(i, v)
	}

	wg.Wait()

	// TOOD Improve this (works for now)
	for i, v := range failedDownloads {
		log.Debug("retrying", "v", v)
		request.DownloadFile(v, failedDirectories[i])
	}

	return true
}
