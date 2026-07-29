package autostarter

import (
	"fmt"
	"nvd/internal/logger"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/emersion/go-autostart"
)

const linuxServiceName string = "nvd"

func AddToAutostart(dl *logger.DownloaderLogger) (ok bool, err error) {
	exePath, err := os.Executable()
	if err != nil {
		return false, err
	}
	exePath, _ = filepath.Abs(exePath)

	app := &autostart.App{
		Name:        "nvd",
		DisplayName: "Nimk`s Video Downloader",
		Exec:        []string{exePath},
	}

	if app.IsEnabled() {
		return false, nil
	} else {
		if err := app.Enable(); err != nil {
			return false, err
		}

		if runtime.GOOS == "windows" {
			return true, nil
		} else {
			dl.LogInfo("Program added to autostart")

			cmd := exec.Command(exePath, "--daemon")

			cmd.Stdout = nil
			cmd.Stderr = nil
			cmd.Stdin = nil

			err := cmd.Start()
			if err != nil {
				return false, fmt.Errorf("backround launching error: %w", err)
			}

			os.Exit(0)
		}
	}

	return true, nil
}
