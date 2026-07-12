package logger

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

type DownloaderLogger struct {
	logger   *slog.Logger
	filePath string
}

func InitDownloaderLogger(filePath string) (*DownloaderLogger, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error("Cannot open log file")
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(file, nil))

	dl := DownloaderLogger{logger: logger, filePath: filePath}

	dl.LogInfo("App started!")

	return &dl, nil
}

func (l *DownloaderLogger) OpenLogFile() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	exeDir := filepath.Dir(exePath)

	fullFilePath := filepath.Join(exeDir, l.filePath)

	cmd := exec.Command("notepad.exe", fullFilePath)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x00000010,
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

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
