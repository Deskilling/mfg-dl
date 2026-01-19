package request

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"mfg-dl/pkg/filesystem"

	"github.com/charmbracelet/log"
)

func DownloadFile(url, filePath string) error {
	err := filesystem.CreatePath(filePath)
	if err != nil {
		err = fmt.Errorf("failed ensuring dir: %w", err)
		log.Error(err)
		return err
	}

	out, err := os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("failed to create file: %w", err)
		log.Error(err)
		return err
	}
	defer out.Close()

	req, err := NewRequest(http.MethodGet, url)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		log.Error(err)
		return err
	}

	resp, err := Client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed request: %w", err)
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
		return err
	}

	if _, err = io.Copy(out, resp.Body); err != nil {
		err = fmt.Errorf("failed to download file: %w", err)
		log.Error(err)
		return err
	}

	return nil
}
