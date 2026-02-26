package filesystem

import (
	"os"
)

func ReadFile(filepath string) (content string, err error) {
	if !ExistPath(filepath) {
		return "", nil
	}

	byte, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(byte), err
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
