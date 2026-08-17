package logger

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const logFileName = "app.log"

type DownloaderLogger struct {
	logger *slog.Logger
	file   *os.File
}

func InitDownloaderLogger() (*DownloaderLogger, error) {
	dl := &DownloaderLogger{}

	dl.rotate()

	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error("Cannot open log file")
		return nil, err
	}
	dl.file = file

	dl.logger = slog.New(slog.NewTextHandler(file, nil))

	dl.LogInfo("App started!")

	return dl, nil
}

func (l *DownloaderLogger) rotate() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("exe path retrieving error: %w", err)
	}

	exeDir := filepath.Dir(exePath)

	oldLogsPath := filepath.Join(exeDir, "old-logs")
	logFilePath := filepath.Join(exeDir, "app.log")

	if err := os.MkdirAll(oldLogsPath, 0755); err != nil {
		return fmt.Errorf("old logs directory creation error: %w", err)
	}

	files, err := os.ReadDir(oldLogsPath)
	if err != nil {
		return fmt.Errorf("old logs directory read error: %w", err)
	}

	type fileInfo struct {
		entry   os.DirEntry
		modTime time.Time
	}
	var filteredFiles []fileInfo

	if err := os.Rename(logFilePath, fmt.Sprintf("%s/app %v.log", oldLogsPath, time.Now().Format("2006-01-02_15-04-05"))); err != nil {
		return fmt.Errorf("log file moving error: %w", err)
	}

	file, err := os.Create(logFilePath)
	if err != nil {
		return fmt.Errorf("log file creation error: %w", err)
	}
	file.Close()

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		info, err := f.Info()
		if err != nil {
			return fmt.Errorf("old log file reading error: %w", err)
		}
		filteredFiles = append(filteredFiles, fileInfo{entry: f, modTime: info.ModTime()})
	}
	filesCount := len(filteredFiles) + 1

	if filesCount <= 3 {
		return nil
	}

	sort.Slice(filteredFiles, func(i, j int) bool {
		return filteredFiles[i].modTime.Before(filteredFiles[j].modTime)
	})

	toRemoveCount := filesCount - 3

	for i := 0; i < toRemoveCount; i++ {
		fileToDelete := filteredFiles[i].entry
		fullPath := filepath.Join(oldLogsPath, fileToDelete.Name())

		err := os.Remove(fullPath)
		if err != nil {
			return fmt.Errorf("old log file removing error: %w", err)
		}
	}

	return nil
}

func (l *DownloaderLogger) LogInfo(s string) {
	l.logger.Info(s)
}

func (l *DownloaderLogger) LogError(s string) {
	l.logger.Error(s)
}

func (l *DownloaderLogger) LogFatal(s string) {
	l.logger.Error(fmt.Sprintf("FATAL: %s", s))
}

func (l *DownloaderLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.LogInfo(fmt.Sprintf("Got %s request on localhost:8080%s from %s", r.Method, r.URL.Path, r.RemoteAddr))
		next.ServeHTTP(w, r)
	})
}
