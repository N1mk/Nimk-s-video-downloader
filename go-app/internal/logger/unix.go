//go:build !windows

package logger

import (
	"os/exec"
)

func setPlatformSysProcAttr(cmd *exec.Cmd) {}
