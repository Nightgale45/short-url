package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	once     sync.Once
	instance *slog.Logger
)

// returns the singleton instance of the logger
func GetInstance() *slog.Logger {
	once.Do(func() {
		instance = slog.New(logHandler())
	})
	return instance
}

func logHandler() slog.Handler {

	option := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}

	return slog.NewTextHandler(os.Stderr, &option)
}
