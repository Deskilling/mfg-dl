package m3u

import (
	"fmt"
	"mfg-dl/filesystem"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/charmbracelet/log"
)

func ConvertTSFilesToVideo(directory, output string) (bool, error) {
	filesystem.EnsureDir(output)

	tsFiles, err := filesystem.GetAllFilesFromDirectory(directory, ".ts")
	if err != nil {
		err = fmt.Errorf("failed to get .ts files: %w", err)
		log.Error(err)
		return false, err
	}

	if len(tsFiles) == 0 {
		err = fmt.Errorf("no .ts files found in directory %s", directory)
		log.Error(err)
		return false, err
	}

	sort.Slice(tsFiles, func(i, j int) bool {
		nameI := tsFiles[i].Name()
		nameJ := tsFiles[j].Name()

		numStrI := nameI[:len(nameI)-len(filepath.Ext(nameI))]
		numStrJ := nameJ[:len(nameJ)-len(filepath.Ext(nameJ))]

		numI, errI := strconv.Atoi(numStrI)
		numJ, errJ := strconv.Atoi(numStrJ)

		if errI != nil || errJ != nil {
			return nameI < nameJ
		}
		return numI < numJ
	})

	listFile, err := os.Create(directory + "/segments.txt")
	if err != nil {
		err = fmt.Errorf("failed to create temporary list file: %w", err)
		log.Error(err)
		return false, err
	}
	defer listFile.Close()

	for _, file := range tsFiles {
		_, err := listFile.WriteString(fmt.Sprintf("file '%s'\n", filepath.Join(directory, file.Name())))
		if err != nil {
			err = fmt.Errorf("failed to write to list file: %w", err)
			log.Error(err)
			return false, err
		}
	}

	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", listFile.Name(),
		"-c", "copy",
		"-nostats",
		"-y",
		output,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Info("Running FFmpeg command", "command", cmd.String())

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("ffmpeg command failed: %w", err)
		log.Error(err)
		return false, err
	}

	return true, nil
}
