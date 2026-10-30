package core

import (
	"log/slog"

	"github.com/pterm/pterm"
)

// i personally dont like the fact that the logger colors are
// dependent on pterm but it just looks better and is simpler to do

func InitLogger(level int) {
	isDebug := level == -4

	pterm.DefaultLogger.ShowCaller = isDebug
	pterm.DefaultLogger.Level = pterm.LogLevel(level)

	handler := pterm.NewSlogHandler(&pterm.DefaultLogger)
	slog.SetDefault(slog.New(handler))
}
