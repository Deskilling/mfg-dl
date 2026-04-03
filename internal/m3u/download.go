package m3u

import (
	"fmt"
	"sync"
	"time"

	"mfg-dl/internal/core"
	"mfg-dl/internal/request"
	"mfg-dl/internal/tui/components"

	"charm.land/log/v2"
)

func DownloadSegments(index Index, baseURL, directory string) (err error) {
	type job struct {
		url  string
		file string
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, core.GetConfig().Downloads.MaxSegmentConcurrency)

	var retryMu sync.Mutex
	var retryQueue []job

	segmentCnt := len(index.Segments)

	for i, seg := range index.Segments {
		wg.Add(1)
		semaphore <- struct{}{}

		components.PrintProgress(i+1, segmentCnt)
		go func(i int, seg Segment) {
			defer wg.Done()
			defer func() { <-semaphore }()

			file := fmt.Sprintf("%s%d.ts", directory, i)
			url := baseURL + seg.URI

			err := request.DownloadFile(url, file)
			if err != nil {
				log.Error("Failed to download", "url", url, "file", file)
				retryMu.Lock()
				retryQueue = append(retryQueue, job{url, file})
				retryMu.Unlock()
			}
		}(i, seg)
	}

	wg.Wait()

	retryDelay := time.Duration(core.GetConfig().Downloads.RetryDelay) * time.Second
	for _, j := range retryQueue {
		for attempt := range core.GetConfig().Downloads.MaxRetires {
			log.Debug("Retry download", "url", j.url, "path", j.file)
			err = request.DownloadFile(j.url, j.file)
			if err != nil {
				if attempt == core.GetConfig().Downloads.MaxRetires-1 {
					return
				}
			} else {
				log.Debug("Redownloaded unc")
				break
			}
			time.Sleep(retryDelay)
		}
	}

	return
}
