package autostarter

import (
	"nvd/internal/logger"
	"nvd/internal/project_errors"
	"runtime"
)

func AddToAutostart(dl *logger.DownloaderLogger) (ok bool, err error) {
	switch runtime.GOOS {
	case "windows":
		return addToAutostartWindows()
	case "linux":
		return addToAutostartUnix(dl)
	case "darwin":
		return addToAutostartUnix(dl)
	default:
		return false, project_errors.ErrUnknownOS
	}
}
