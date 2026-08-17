//go:build !windows

package loader

import (
	"os"
	"os/exec"
	"path/filepath"
)

const ytDlpName = "yt-dlp"

func (d *Loader) UpdatePath() error {

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Dir(exePath)
	d.ytdlpPath = filepath.Join(dir, ytDlpName)

	return nil
}

func setPlatformSysProcAttr(cmd *exec.Cmd) {}
