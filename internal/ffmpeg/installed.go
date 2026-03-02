package ffmpeg

import (
	"os/exec"

	"github.com/charmbracelet/log"
)

func CheckInstalled() (err error) {
	cmd := exec.Command("ffmpeg", "-version")
	err = cmd.Run()

	if err != nil {
		log.Info("ffmpeg is not installed")
	}

	return err
}
