package request

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"mfg-dl/filesystem"

	"github.com/charmbracelet/log"
)

func DownloadFile(url, filePath string) error {
	err := filesystem.EnsureDir(filePath)
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

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		log.Error(err)
		return err
	}
	req.Header.Set("User-Agent", "Deskilling/aniworld-dl")

	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed request: %w", err)
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
		log.Error(err)
		return err
	}

	if _, err = io.Copy(out, resp.Body); err != nil {
		err = fmt.Errorf("failed to download file: %w", err)
		log.Error(err)
		return err
	}

	return nil
}
