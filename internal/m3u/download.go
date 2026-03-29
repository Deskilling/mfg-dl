package m3u

import (
	"fmt"
	"sync"
	"time"

	"mfg-dl/internal/core"
	"mfg-dl/internal/request"
)

func DownloadSegments(index Index, baseURL, directory string) error {
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
				retryMu.Lock()
				retryQueue = append(retryQueue, job{url, file})
				retryMu.Unlock()
			} else {
				// TOOD
			}
		}(i, seg)
	}

	wg.Wait()

	retryDelay := time.Duration(core.GetConfig().Downloads.RetryDelay) * time.Second
	for _, j := range retryQueue {
		for attempt := 0; attempt < core.GetConfig().Downloads.MaxRetires; attempt++ {
			err := request.DownloadFile(j.url, j.file)
			if err == nil {
				break
			}
			time.Sleep(retryDelay)
		}
	}

	return nil
}
