package core

import (
	"log/slog"
	"os"
	"time"

	"github.com/Deskilling/gopkg/pkg/filesystem"
)

func CleanUp() {
	files, err := filesystem.ReadDirectoryRecursive(GetConfig().Location.Temp+"/cache", "")
	if err != nil {
		return
	}
	for _, file := range files {
		stat, err := os.Stat(file)
		if err != nil {
			continue
		}

		if time.Since(stat.ModTime()) > time.Minute*time.Duration(GetConfig().Cache.Minutes) {
			err = filesystem.DeleteFile(file)
			if err != nil {
				return
			}
		}
	}
}

func CreateCleanupTask() {
	slog.Info("Starting Cleanup Task")
	go func() {
		CleanUp()
		ticker := time.NewTicker(time.Minute * time.Duration(GetConfig().Cache.CleanMinutes))
		defer ticker.Stop()
		for range ticker.C {
			CleanUp()
		}
	}()
}
