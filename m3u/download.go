package m3u

import (
	"mfg-dl/request"
	"strconv"
	"sync"

	"github.com/charmbracelet/log"
)

var maxConcurrency = 64

func DownloadSegments(index *Index, baseURL, directory string) bool {
	var (
		failedDownloads []string
		failedFiles     []string
	)

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	semaphore := make(chan struct{}, maxConcurrency)
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
				mu.Lock()
				failedDownloads = append(failedDownloads, baseURL+v.URI)
				failedFiles = append(failedFiles, directory+s+".ts")
				mu.Unlock()
				log.Error(err)
			}
		}(i, v)
	}

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if len(failedDownloads) != len(failedFiles) {
		log.Debug(failedFiles)
		log.Debug(failedDownloads)
		log.Fatal("Mismatch")
	}

	var done []int
	for i, v := range failedDownloads {
		log.Debug("Retrying", "file", v)
		err := request.DownloadFile(v, failedFiles[i])
		if err == nil {
			done = append(done, i)
		}
	}

	for i := len(done) - 1; i >= 0; i-- {
		idxToRemove := done[i]
		failedDownloads = append(failedDownloads[:idxToRemove], failedDownloads[idxToRemove+1:]...)
		failedFiles = append(failedFiles[:idxToRemove], failedFiles[idxToRemove+1:]...)
	}

	if len(failedDownloads) != 0 && len(failedFiles) != 0 {
		for i, v := range failedDownloads {
			log.Debug("Retrying", "file", v)
			err := request.DownloadFile(v, failedFiles[i])
			if err != nil {
				return false
			}
		}
	}

	return true
}
