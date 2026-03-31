package ffmpeg

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"

	"mfg-dl/internal/tui/components"
)

func CheckInstalled() (err error) {
	cmd := exec.Command("ffmpeg", "-version")
	err = cmd.Run()

	if err != nil {
		if runtime.GOOS == "windows" {
			input, err := components.ReadString(components.Reader, "y/n: ")
			if err != nil {
				return fmt.Errorf("failed userinput: %w", err)
			}

			if input == "y" {
				cmd = exec.Command("winget", "install", "Gyan.FFmpeg")
				err = cmd.Run()

				if err != nil {
					return errors.New("Failed to install using Winget")
				}
			}
		} else {
			return errors.New("FFmpeg is not installed")
		}
	}

	return nil
}
