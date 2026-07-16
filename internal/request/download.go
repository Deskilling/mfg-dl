package request

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"charm.land/log/v2"
	"github.com/Deskilling/gopkg/pkg/filesystem"
)

func DownloadFile(client *http.Client, url, filePath string) (err error) {
	if url == "" {
		log.Error("invalid url to request")
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
	defer out.Close()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", userAgent)
	resp, err := Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	return err
}
