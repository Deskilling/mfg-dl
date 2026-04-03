package request

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"mfg-dl/pkg/filesystem"
)

func DownloadFile(url, filePath string) (err error) {
	err = filesystem.CreatePath(filePath)
	if err != nil {
		return
	}

	out, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer out.Close()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", userAgent)
	resp, err := Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	return
}
