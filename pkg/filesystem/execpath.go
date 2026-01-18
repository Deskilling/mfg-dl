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
