package filesystem

import (
	"errors"
	"os"
	"path/filepath"
)

var execPath string

func InitExecDir() (err error) {
	execPath, err = os.Executable()
	if err != nil {
		return err
	}

	execPath = filepath.Dir(execPath)
	if execPath == "" {
		return errors.New("failed to get directory of executable")
	}
	return nil
}

func GetExecPath() *string {
	return &execPath
}

func ChangeExecPath() (err error) {
	if execPath == "" {
		err = InitExecDir()
		if err != nil {
			return err
		}
	}

	err = os.Chdir(execPath)
	if err != nil {
		return err
	}

	return nil
}
