package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/Deskilling/gopkg/pkg/filesystem"
)

func ConvertTSFilesToVideo(directory, output string, args []string) (err error) {
	filesystem.CreatePath(output)

	tsFiles, err := filesystem.ReadDirectory(directory, ".ts")
	if err != nil {
		return err
	}

	if len(tsFiles) == 0 {
		return fmt.Errorf("no .ts files found in directory %s", directory)

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

	listFile, err := os.Create(directory + "segments.txt")
	if err != nil {
		return err
	}
	defer listFile.Close()

	for _, file := range tsFiles {
		_, err := fmt.Fprintf(listFile, "file '%s'\n", file.Name())
		if err != nil {
			return err
		}
	}

	var finalArgs []string
	for _, arg := range args {
		switch arg {
		case "input":
			finalArgs = append(finalArgs, "-i", listFile.Name())
		case "output":
			finalArgs = append(finalArgs, output)
		default:
			finalArgs = append(finalArgs, arg)
		}
	}

	cmd := exec.Command("ffmpeg", finalArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
