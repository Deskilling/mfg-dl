package core

import (
	"mfg-dl/pkg/filesystem"
	"time"

	"charm.land/log/v2"
)

func CleanUp() {
	files, err := filesystem.ReadDirectoryRecursive(GetConfig().Location.Cache+"/", "")
	if err != nil {
		log.Errorf("failed to read cache directory: %v", err)
		return
	}
	for _, file := range files {
		stat, err := file.Info()
		if err == nil && time.Since(stat.ModTime()) > time.Minute*time.Duration(GetConfig().Cache.Minutes) {
			filesystem.DeleteFile(file.Name())
		}
	}
}

func CreateCleanupTask() {
	go func() {
		ticker := time.NewTicker(time.Hour * time.Duration(GetConfig().Cache.CleanMinutes))
		defer ticker.Stop()
		for range ticker.C {
			CleanUp()
		}
	}()
}
