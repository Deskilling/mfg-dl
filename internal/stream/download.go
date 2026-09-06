package stream

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"mfg-dl/internal/core"
	"mfg-dl/internal/request"

	"charm.land/log/v2"
	"github.com/Deskilling/gopkg/pkg/m3u8"
	"github.com/cavaliergopher/grab/v3"
)

var Client = grab.NewClient()

func init() {
	Client.HTTPClient = &http.Client{
		Timeout:   60 * time.Second,
		Transport: request.Client.Transport,
	}
}

func DownloadSegments(index m3u8.Index, baseURL, directory string) (err error) {
	type job struct {
		url  string
		file string
	}

	cfg := core.GetConfig().Downloads

	reqs := make([]*grab.Request, 0, len(index.Segments))
	for i, seg := range index.Segments {
		file := fmt.Sprintf("%s%d.ts", directory, i)
		url := baseURL + seg.URI

		req, err := grab.NewRequest(file, url)
		if err != nil {
			return fmt.Errorf("building request for segment %d: %w", i, err)
		}
		reqs = append(reqs, req)
	}

	respCh := Client.DoBatch(cfg.MaxSegmentConcurrency, reqs...)

	var failed []job
	for resp := range respCh {
		if err := resp.Err(); err != nil {
			log.Error("Failed to download", "url", resp.Request.URL(), "file", resp.Filename, "err", err)
			failed = append(failed, job{url: resp.Request.URL().String(), file: resp.Filename})
		}
	}

	retryDelay := time.Duration(cfg.RetryDelay) * time.Second
	var stillFailed []error
	for _, j := range failed {
		var lastErr error
		for attempt := range cfg.MaxRetires {
			req, err := grab.NewRequest(j.file, j.url)
			if err != nil {
				lastErr = err
				break
			}
			resp := Client.Do(req)
			if lastErr = resp.Err(); lastErr == nil {
				break
			}
			if attempt < cfg.MaxRetires-1 {
				time.Sleep(retryDelay)
			}
		}
		if lastErr != nil {
			log.Error("GG on segment", "url", j.url, "path", j.file, "err", lastErr)
			stillFailed = append(stillFailed, fmt.Errorf("%s: %w", j.file, lastErr))
		}
	}

	if len(stillFailed) > 0 {
		return fmt.Errorf("%d segment(s) failed permanently: %w", len(stillFailed), errors.Join(stillFailed...))
	}

	return nil
}
