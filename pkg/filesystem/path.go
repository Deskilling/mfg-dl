package filesystem

import (
	"os"
	"path/filepath"
	"strings"
)

func ExistPath(path string) (exists bool) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return true
	}
	return false
}

func GetSlug(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

func BackOncePath(path, after string) (newPath string) {
	return strings.TrimRight(path, after)
}

func CleanPath(path string) (cleanedPath string) {
	path = strings.TrimRight(path, `/\`)
	path = filepath.Clean(path)
	sep := string(filepath.Separator)
	if !strings.HasSuffix(path, sep) {
		path += sep
	}

	return path
}

func CreatePath(path string) (err error) {
	if ExistPath(path) {
		return nil
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	return nil
}

func ClearPath(path, extension string) (err error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		info, err := entry.Info()
		if err != nil {
			return err
		}

		if info.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != extension {
			continue
		}

		if err := os.Remove(fullPath); err != nil {
			return err
		}
	}

	return nil
}

func IsDirEmpty(path string) (empty bool) {
	if !ExistPath(path) {
		return true
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return true
	}

	if len(files) != 0 {
		return false
	}

	return true
}
