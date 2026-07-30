//go:build !windows

package convertor

import (
	"os/exec"
)

func setPlatformSysProcAttr(cmd *exec.Cmd) {}
