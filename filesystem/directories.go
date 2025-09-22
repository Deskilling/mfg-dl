package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func EnsureDir(path string) error {
	dirPath := path
	if filepath.Ext(path) != "" {
		dirPath = filepath.Dir(path)
	}

	log.Debug("creating filepath", "path", path)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	return nil
}

/*
func checkStringValidPath(path string) string {
	if len(path) == 0 {
		return ""
	}

	lastChar := path[len(path)-1:]
	if lastChar != string(filepath.Separator) {
		path += string(filepath.Separator)
	}

	_, err := DoesPathExist(path)
	if err != nil {
		return ""
	}
	return path
}
*/
