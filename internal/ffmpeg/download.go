package ffmpeg

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"mfg-dl/pkg/filesystem"

	"charm.land/log/v2"
)

func GetTotalMicro(url string) (totalMicro int, err error) {
	args := []string{
		"-v", "error",
		"-i", url,
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
	}

	out, err := exec.Command("ffprobe", args...).Output()
	if err != nil {
		return 0, err
	}

	durStr := strings.TrimSpace(string(out))
	if durStr == "" {
		return 0, nil
	}

	durSec, err := strconv.ParseFloat(durStr, 64)
	if err != nil {
		return 0, err
	}

	return int(durSec * 1000 * 1000), nil
}

func DownloadHLS(url, output string) (err error) {
	filesystem.CreatePath(output)

	totalMicro, _ := GetTotalMicro(url)
	log.Debug("Got toalms", "totalm", totalMicro)

	args := []string{
		"-y",
		"-user_agent", "Mozilla/5.0",
		"-reconnect", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "5",
		"-progress", "pipe:1",
		"-i", url,
		"-c", "copy",
		output,
	}

	cmd := exec.Command("ffmpeg", args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	cmd.Stderr = nil

	err = cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			_, ok := strings.CutPrefix(line, "out_time_us=")
			if ok {
				log.Infof("Unc %s/%v", line, totalMicro)
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return err
	}

	fmt.Println("Download complete.")
	return nil
}
