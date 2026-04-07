package filesystem

import (
	"errors"
	"os"
)

func ReadFile(filepath string) (content []byte, err error) {
	if !ExistPath(filepath) {
		return nil, errors.New("file does not exist")
	}

	byte, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return byte, err
}

func WriteFile(path string, byte []byte) (err error) {
	if !ExistPath(path) {
		CreatePath(path)
	}

	err = os.WriteFile(path, byte, 0777)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFile(path string) (err error) {
	if ExistPath(path) {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}
