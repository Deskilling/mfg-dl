package ffmpeg

import (
	"os/exec"

	"charm.land/log/v2"
)

func CheckInstalled() (err error) {
	cmd := exec.Command("ffmpeg", "-version")
	err = cmd.Run()

	if err != nil {
		log.Info("ffmpeg is not installed")
	}

	return err
}
