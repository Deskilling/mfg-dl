package ffmpeg

import (
	"os/exec"
	"runtime"

	"mfg-dl/internal/tui/components"

	"charm.land/log/v2"
)

func CheckInstalled() (err error) {
	cmd := exec.Command("ffmpeg", "-version")
	err = cmd.Run()

	if err != nil {
		if runtime.GOOS == "windows" {
			log.Info("FFmpeg is not installed\nDo you want to install it using winget")
			input := components.ReadString(components.Reader, "y/n: ")

			if input == "y" {

				cmd = exec.Command("winget", "install", "Gyan.FFmpeg")
				err = cmd.Run()

				if err != nil {
					log.Error("Winget is not installed or error happend :(")
				}
			}
		} else {
			log.Info("ffmpeg is not installed")
		}
	}

	return err
}
