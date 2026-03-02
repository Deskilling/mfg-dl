package filesystem

import (
	"os"
	"path/filepath"
)

var execDir string

func InitExecDir() (err error) {
	execDir, err = os.Executable()
	if err != nil {
		return err
	}

	execDir = filepath.Dir(execDir)
	return nil
}

func GetExecDir() *string {
	return &execDir
}

func ChangeExecDir() (err error) {
	if execDir == "" {
		err = InitExecDir()
		if err != nil {
			return err
		}
	}

	err = os.Chdir(execDir)
	if err != nil {
		return err
	}

	return nil
}
