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
		switch runtime.GOOS {
		case ("windows"):
			fmt.Println("Do you want to install ffmpeg (Gyan.FFmpeg) using winget?")
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

				return nil
			}

		case ("darwin"):
			fmt.Println("Do you want to install ffmpeg using homebrew?")
			input, err := components.ReadString(components.Reader, "y/n: ")
			if err != nil {
				return fmt.Errorf("failed userinput: %w", err)
			}

			if input == "y" {
				cmd = exec.Command("brew", "install", "ffmpeg")
				err = cmd.Run()

				if err != nil {
					return errors.New("Failed to install using Winget")
				}

				return nil
			}
		}

		return errors.New("FFmpeg is not installed")
	}

	return nil
}
