package ffmpeg

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
)

func CheckInstalled() (err error) {
	cmd := exec.Command("ffmpeg", "-version")
	err = cmd.Run()

	if err != nil {
		switch runtime.GOOS {
		case ("windows"):
			fmt.Println("Do you want to install ffmpeg (Gyan.FFmpeg) using winget?")
			// TODO

			/*
				input, err := components.ReadString("y/n: ")
				if err != nil {
					return fmt.Errorf("failed userinput: %w", err)
				}
			*/

			input := "n"

			if input == "y" {
				cmd = exec.Command("winget", "install", "Gyan.FFmpeg")
				err = cmd.Run()

				if err != nil {
					return errors.New("failed to install using Winget")
				}

				return nil
			}

		case ("darwin"):
			fmt.Println("Do you want to install ffmpeg using homebrew?")
			// TODO

			/*
				input, err := components.ReadString("y/n: ")
				if err != nil {
					return fmt.Errorf("failed userinput: %w", err)
				}
			*/

			input := "n"

			if input == "y" {
				cmd = exec.Command("brew", "install", "ffmpeg")
				err = cmd.Run()

				if err != nil {
					return errors.New("failed to install using Winget")
				}

				return nil
			}
		}

		return errors.New("ffmpeg is not installed")
	}

	return nil
}
