//go:build !windows

package convertor

import (
	"os"
	"os/exec"
	"path/filepath"
)

const ffmpegName = "ffmpeg"

func (c *Convertor) UpdatePath() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Dir(exePath)
	c.ffmpegPath = filepath.Join(dir, ffmpegName)

	return nil
}

func setPlatformSysProcAttr(cmd *exec.Cmd) {}
