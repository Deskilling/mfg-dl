package m3u

import (
	"fmt"
	"sync"
	"time"

	"mfg-dl/internal/core"
	"mfg-dl/internal/request"

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

	for i, seg := range index.Segments {
		wg.Add(1)
		semaphore <- struct{}{}

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
			err := request.DownloadFile(j.url, j.file)
			if err == nil {
				if attempt == core.GetConfig().Downloads.MaxRetires {
					log.Error("Failed downloading segment", "attempt", attempt, "url", j.url, "file", j.file, "err", err)
					return err
				}
			}
			time.Sleep(retryDelay)
		}
	}

	return nil
}
