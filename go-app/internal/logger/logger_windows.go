//go:build windows

package logger

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const notepadAppName = "notepad.exe"

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

func setPlatformSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x00000010,
	}
}
