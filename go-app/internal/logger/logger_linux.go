//go:build linux

package logger

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const notepadAppName = "xdg-open"

func (l *DownloaderLogger) OpenLogFile() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)
	fullFilePath := filepath.Join(exeDir, logFileName)

	cmd := exec.Command(notepadAppName, fullFilePath)

	setPlatformSysProcAttr(cmd)

	if err := cmd.Start(); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	return nil
}

func setPlatformSysProcAttr(cmd *exec.Cmd) {}
