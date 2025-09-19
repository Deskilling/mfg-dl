package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

var (
	execDir string
	err     error
)

func InitExecDir() {
	execDir, err = os.Executable()
	if err != nil {
		err = fmt.Errorf("failed to get executable path: %w", err)
		log.Error(err)
		os.Exit(1)
	}

	execDir = filepath.Dir(execDir)
	log.Debug("execDir", "path", execDir)
}

func GetExecDir() string {
	return execDir
}
