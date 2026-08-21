//go:build linux

package loader

import (
	"os"
	"os/exec"
	"path/filepath"
)

const (
	ytDlpName = "yt-dlp"
	qjsName   = "qjs"
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

func setPlatformSysProcAttr(cmd *exec.Cmd) {}
