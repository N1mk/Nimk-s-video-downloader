package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type DownloaderLogger struct {
	logger *slog.Logger
}

func InitDownloaderLogger(fileName string) (*DownloaderLogger, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error("Cannot open log file")
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(file, nil))

	dl := DownloaderLogger{logger: logger}

	dl.LogInfo("App started!")

	return &dl, nil
}

func (l *DownloaderLogger) LogInfo(s string) {
	l.logger.Info(s)
}

func (l *DownloaderLogger) LogError(s string) {
	l.logger.Error(s)
}

func (l *DownloaderLogger) LogFatal(s string) {
	l.logger.Error(fmt.Sprintf("FATAL ERROR: %s", s))
}
