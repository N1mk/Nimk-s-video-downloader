package autostarter

import (
	"os"

	"golang.org/x/sys/windows/registry"
)

func AddToAutostart() (bool, error) {
	exePath, err := os.Executable()
	if err != nil {
		return false, err
	}

	k, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return false, err
	}
	defer k.Close()

	val, _, err := k.GetStringValue("VideoDownloader")
	if err == nil && exePath == val {
		return false, nil
	}

	err = k.SetStringValue("VideoDownloader", exePath)
	if err != nil {
		return false, err
	}

	return true, nil
}
