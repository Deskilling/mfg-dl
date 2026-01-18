package util

import "github.com/charmbracelet/log"

func InitLogger(level log.Level) {
	log.SetLevel(level)

	if level == log.DebugLevel {
		log.SetReportCaller(true)
	} else {
		log.SetReportCaller(false)
		log.SetReportTimestamp(false)
	}

	log.SetPrefix("mfg-dl")
}
