package ffmpeg

import (
	"os/exec"
	"runtime"

	"charm.land/log/v2"
)

func CheckInstalled() (err error) {
	cmd := exec.Command("ffmpeg", "-version")
	err = cmd.Run()

	if err != nil {
		if runtime.GOOS == "windows" {
			log.Info("Downloading ffmpeg with winget")

			cmd = exec.Command("winget", "install", "Gyan.FFmpeg")
			err = cmd.Run()

			if err != nil {
				log.Error("Winget is not installed or error happend :(")
			}
		} else {
			log.Info("ffmpeg is not installed")
		}
	}

	return err
}
