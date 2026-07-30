//go:build !windows

package downloader

import (
	"os/exec"
)

func setPlatformSysProcAttr(cmd *exec.Cmd) {}
