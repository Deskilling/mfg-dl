package m3u

// TODO Remove request dependency
import (
	"fmt"
	"strconv"
	"sync"

	"mfg-dl/internal/request"
)

// TODO test for best value
// TODO irgend was kochen (maybe average check oder so)
var maxConcurrency = 16

// TODO Add multithread download for failed files
func DownloadSegments(index Index, baseURL, directory string) (err error) {
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

			s := strconv.Itoa(i)
			err := request.DownloadFile(baseURL+v.URI, directory+s+".ts")
			if err != nil {
				mu.Lock()
				failedDownloads = append(failedDownloads, baseURL+v.URI)
				failedFiles = append(failedFiles, directory+s+".ts")
				mu.Unlock()
			}
		}(i, v)
	}

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if len(failedDownloads) != len(failedFiles) {
		return fmt.Errorf("failed downloads and files missmatch")
	}

	var done []int
	for i, v := range failedDownloads {
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
			err := request.DownloadFile(v, failedFiles[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
