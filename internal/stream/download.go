package stream

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"mfg-dl/internal/core"
	"mfg-dl/internal/request"

	"github.com/Deskilling/gopkg/pkg/m3u8"
	"github.com/cavaliergopher/grab/v3"
	"github.com/pterm/pterm"
)

var Client = grab.NewClient()

func init() {
	Client.HTTPClient = &http.Client{
		Timeout: 60 * time.Second,
		// TODO Look into better transport settings
		Transport: request.Client.Transport,
	}
}

func DownloadSegments(index m3u8.Index, baseURL, directory string) (err error) {
	type job struct {
		url  string
		file string
		err  error
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

	p, _ := pterm.DefaultProgressbar.
		WithTotal(len(reqs)).
		WithTitle("Downloading Segments").
		WithRemoveWhenDone(true).
		Start()

	var failed []job
	for resp := range respCh {
		err := resp.Err()
		p.Increment()

		if err != nil {
			slog.Error("Failed to download", "url", resp.Request.URL(), "file", resp.Filename, "err", err)
			failed = append(failed, job{
				url:  resp.Request.URL().String(),
				file: resp.Filename,
				err:  err,
			})
		}
	}

	retryDelay := time.Duration(cfg.RetryDelay) * time.Second
	var stillFailed []error
	for attempt := range cfg.MaxRetires {
		retryReqs := make([]*grab.Request, 0, len(failed))
		for _, j := range failed {
			req, err := grab.NewRequest(j.file, j.url)
			if err != nil {
				stillFailed = append(stillFailed, fmt.Errorf("%s: %w", j.file, err))
				continue
			}
			retryReqs = append(retryReqs, req)
		}

		failed = nil
		for resp := range Client.DoBatch(cfg.MaxSegmentConcurrency, retryReqs...) {
			err := resp.Err()
			if err != nil {
				failed = append(failed, job{
					url:  resp.Request.URL().String(),
					file: resp.Filename,
					err:  err,
				})
			}
		}

		if len(failed) > 0 && attempt < cfg.MaxRetires-1 {
			time.Sleep(retryDelay)
		}
	}

	for _, j := range failed {
		slog.Error("GG on segment", "url", j.url, "path", j.file, "err", j.err)
		stillFailed = append(stillFailed, fmt.Errorf("%s: %w", j.file, j.err))
	}

	if len(stillFailed) > 0 {
		return fmt.Errorf("%d segment(s) failed permanently: %w", len(stillFailed), errors.Join(stillFailed...))
	}

	return nil
}
