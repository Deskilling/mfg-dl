package core

import (
	"mfg-dl/pkg/filesystem"
	"os"
	"time"

	"charm.land/log/v2"
)

func CleanUp() {
	files, err := filesystem.ReadDirectoryRecursive(GetConfig().Location.Temp+"/cache", "")
	if err != nil {
		log.Errorf("failed to read cache directory: %v", err)
		return
	}
	for _, file := range files {
		stat, _ := os.Stat(file)
		if time.Since(stat.ModTime()) > time.Minute*time.Duration(GetConfig().Cache.Minutes) {
			err = filesystem.DeleteFile(file)
			if err != nil {
				log.Errorf("failed to delete file %s: %v", file, err)
				return
			}
		}
	}
}

func CreateCleanupTask() {
	go func() {
		CleanUp()
		ticker := time.NewTicker(time.Hour * time.Duration(GetConfig().Cache.CleanMinutes))
		defer ticker.Stop()
		for range ticker.C {
			CleanUp()
		}
	}()
}
