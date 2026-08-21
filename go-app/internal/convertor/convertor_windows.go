//go:build windows

package convertor

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const ffmpegName = "ffmpeg.exe"

func (c *DefaultConvertor) updatePath() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Dir(exePath)
	c.ffmpegPath = filepath.Join(dir, ffmpegName)

	return nil
}

func setPlatformSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x08000000,
	}
}
