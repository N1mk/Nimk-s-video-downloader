//go:build darwin

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

func (d *Loader) UpdatePath() error {

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Dir(exePath)
	d.ytdlpPath = filepath.Join(dir, ytDlpName)
	d.qjsPath = filepath.Join(dir, qjsName)

	return nil
}

func setPlatformSysProcAttr(cmd *exec.Cmd) {}
