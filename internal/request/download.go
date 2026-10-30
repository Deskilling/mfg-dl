package request

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/Deskilling/gopkg/pkg/filesystem"
)

func DownloadFile(client *http.Client, url, filePath string) (err error) {
	if url == "" {
		slog.Error("invalid url to request")
		return fmt.Errorf("invalid url")
	}

	err = filesystem.CreatePath(filePath)
	if err != nil {
		return err
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", userAgent)
	resp, err := Client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	return err
}
