//go:build windows

package loader

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const (
	ytDlpName = "yt-dlp.exe"
	qjsName   = "qjs.exe"
)

func (l *DefaultLoader) updatePath() error {

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Dir(exePath)
	l.ytdlpPath = filepath.Join(dir, ytDlpName)
	l.qjsPath = filepath.Join(dir, qjsName)

	return nil
}

func setPlatformSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x08000000,
	}
}
