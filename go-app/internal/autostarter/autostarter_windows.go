//go:build windows

package autostarter

import (
	"nvd/internal/logger"
	"os"
	"path/filepath"

	"github.com/emersion/go-autostart"
	lnk "github.com/parsiya/golnk"
	"golang.org/x/sys/windows"
)

func AddToAutostart(_ *logger.DownloaderLogger) (ok bool, err error) {
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
		shortcutDirPath, err := windows.KnownFolderPath(windows.FOLDERID_Startup, 0)

		shotrcut, err := lnk.File(filepath.Join(shortcutDirPath, "nvd.lnk"))
		if err != nil {
			return false, err
		}

		if shotrcut.LinkInfo.LocalBasePath == exePath {
			return false, nil
		} else {
			if err := app.Disable(); err != nil {
				return false, err
			}
			if err := app.Enable(); err != nil {
				return false, err
			}
		}
	} else if err := app.Enable(); err != nil {
		return false, err
	}

	return true, nil
}
