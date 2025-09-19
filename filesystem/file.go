package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadFile(filepath string) string {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return "" //fmt.Errorf("failed to read file %s: %w", filepath, err)
	}
	return string(fileContent)
}

// TODO Create CreateFile Fuction
func CreateFile() {

}

func GetAllFilesFromDirectory(directory string, extension string) ([]os.DirEntry, error) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %s does not exist", directory)
	}

	allFiles, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", directory, err)
	}

	var filteredFiles []os.DirEntry

	for _, file := range allFiles {
		if filepath.Ext(file.Name()) == extension {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles, nil
}
