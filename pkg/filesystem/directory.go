package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
)

func ReadDirectory(path string, extension string) (files []os.DirEntry, err error) {
	if !ExistPath(path) {
		return []os.DirEntry{}, nil
	}

	if IsDirEmpty(path) {
		return []os.DirEntry{}, nil
	}

	allFiles, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range allFiles {
		if extension == "" || filepath.Ext(file.Name()) == extension {
			files = append(files, file)
		}
	}
	return files, nil
}

func ReadDirectoryRecursive(path string, extension string) (files []string, err error) {
	if !ExistPath(path) {
		return []string{}, nil
	}
	if IsDirEmpty(path) {
		return []string{}, nil
	}
	err = filepath.WalkDir(path, func(entry string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if extension == "" || filepath.Ext(d.Name()) == extension {
			files = append(files, entry)
		}
		return nil
	})

	return files, err
}

func CopyDirectory(source string, target string) (err error) {
	err = CreatePath(target)
	if err != nil {
		return err
	}

	files, err := ReadDirectory(source, "")
	if err != nil {
		return err
	}

	for _, f := range files {
		sourcePath := filepath.Join(source, f.Name())
		targetPath := filepath.Join(target, f.Name())

		if f.IsDir() {
			err = CopyDirectory(sourcePath, targetPath)
			if err != nil {
				return err
			}

		} else {
			content, err := os.ReadFile(sourcePath)
			if err != nil {
				return err
			}

			err = os.WriteFile(targetPath, content, os.ModeAppend)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func RemoveDirectory(path string) (err error) {
	if !ExistPath(path) {
		return nil
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		currentPath := filepath.Join(path, f.Name())

		if f.IsDir() {
			err = RemoveDirectory(currentPath)
			if err != nil {
				return err
			}
		} else {
			err = os.Remove(currentPath)
			if err != nil {
				return err
			}
		}
	}

	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}
