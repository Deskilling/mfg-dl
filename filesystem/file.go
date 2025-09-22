package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func ReadFile(filepath string) (string, error) {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		err = fmt.Errorf("failed to read file %s: %w", filepath, err)
		log.Error(err)
		return "", err
	}
	return string(fileContent), nil
}

// TODO Create CreateFile Fuction
func CreateFile() {
}

func GetAllFilesFromDirectory(directory string, extension string) ([]os.DirEntry, error) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = fmt.Errorf("directory %s does not exist", directory)
		log.Error(err)
		return nil, err
	}

	allFiles, err := os.ReadDir(directory)
	if err != nil {
		err = fmt.Errorf("failed to read directory %s: %w", directory, err)
		log.Error(err)
		return nil, err
	}

	var filteredFiles []os.DirEntry

	for _, file := range allFiles {
		if filepath.Ext(file.Name()) == extension {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles, nil
}
