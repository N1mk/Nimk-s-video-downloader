//go:build windows

package convertor

import (
	"os/exec"
	"syscall"
)

func setPlatformSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x08000000,
	}
}
