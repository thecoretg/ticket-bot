package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func SetLogger(verbose, debug, toFile bool, logFilePath string) error {
	level := slog.Level(1000)
	if verbose || debug {
		level = slog.LevelInfo
		if debug {
			level = slog.LevelDebug
		}
	}

	if toFile {
		if logFilePath == "" {
			logFilePath = "ticketbot.log"
		}

		fh, err := newFileHandler(logFilePath, level)
		if err != nil {
			return fmt.Errorf("creating file handler: %w", err)
		}
		slog.SetDefault(slog.New(fh))
		return nil
	}

	handler := newStdoutHandler(level)
	slog.SetDefault(slog.New(handler))
	return nil
}

func newFileHandler(filePath string, level slog.Level) (*slog.JSONHandler, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("opening log file %s: %w", filePath, err)
	}

	return slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: level,
	}), nil
}

func newStdoutHandler(level slog.Level) *slog.TextHandler {
	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
}
