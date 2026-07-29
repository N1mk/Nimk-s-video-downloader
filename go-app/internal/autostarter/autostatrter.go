package autostarter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/windows/registry"
)

const linuxServiceName string = "nvd"

// Зависит от ОС
func AddToAutostart() (ok bool, err error) {
	exePath, err := os.Executable()
	if err != nil {
		return false, err
	}
	exePath, _ = filepath.Abs(exePath)

	switch runtime.GOOS {
	case "windows":
		return AddToAutostartWindows(exePath)
	case "linux":
		return AddToAutostartLinux(exePath)
	case "darwin":
		return AddToAutostartMacOS(exePath)
	}

	return false, fmt.Errorf("unknown OS")
}

func AddToAutostartWindows(exePath string) (ok bool, err error) {
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

func AddToAutostartLinux(exePath string) (ok bool, err error) { // ПРИ ЭТОМ СЦЕНАРРИ ПРОГА ДОЛЖНА БЫТЬ ЗАПУЩЕНА ЧЕРЕЗ sudo С ФЛАГОМ --install, ДОБАВИТЬ ЭТО В README ПЕРЕД ПУШЕМ
	if len(os.Args) > 1 && os.Args[1] == "--install" {
		servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", linuxServiceName)

		if os.Getuid() != 0 {
			return false, fmt.Errorf("root access required (use sudo)")
		}

		serviceConfig := fmt.Sprintf(`[Unit]
		Description=Video and audio downloader for YouTube and other platforms
		After=network.target

		[Service]
		Type=simple
		ExecStart=%s
		Restart=on-failure
		User=root

		[Install]
		WantedBy=multi-user.target
		`, exePath)

		err := os.WriteFile(servicePath, []byte(serviceConfig), 0644)
		if err != nil {
			return false, fmt.Errorf("service file writing error: %w", err)
		}

		if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
			return false, fmt.Errorf("daemon-reload error: %w", err)
		}

		if err := exec.Command("systemctl", "enable", linuxServiceName).Run(); err != nil {
			return false, fmt.Errorf("service enabling error: %w", err)
		}

		if err := exec.Command("systemctl", "start", linuxServiceName).Run(); err != nil {
			return false, fmt.Errorf("service starting error: %w", err)
		}
		os.Exit(0)
	} else {
		return false, nil
	}
	return true, nil
}

func AddToAutostartMacOS(exePath string) (ok bool, err error) {
	return true, nil
}
