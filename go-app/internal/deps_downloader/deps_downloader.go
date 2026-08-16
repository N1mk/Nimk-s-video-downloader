package deps_downloader

import (
	"context"
	"fmt"
	"nvd/internal/logger"
	"os"
	"path/filepath"
)

func DownloadDeps(ctx context.Context, dl *logger.DownloaderLogger) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("exe file location determining error: %s", err.Error())
	}

	exeDir := filepath.Dir(exePath)

	if _, err := os.Stat(filepath.Join(exeDir, "yt-dlp.exe")); os.IsNotExist(err) {
		dl.LogInfo("Yt-dlp not found. Downloading...")
		if err := installYtDlp(ctx, exeDir); err != nil {
			return fmt.Errorf("yt-dlp download error: %s", err.Error())
		}
		dl.LogInfo("Yt-dlp downloaded successfully!")
	} else if err != nil {
		return fmt.Errorf("yt-dlp exist check error: %s", err.Error())
	}

	if _, err := os.Stat(filepath.Join(exeDir, "ffmpeg.exe")); os.IsNotExist(err) {
		dl.LogInfo("FFmpeg not found. Downloading...")
		if err := installFFmpeg(ctx, exeDir); err != nil {
			return fmt.Errorf("ffmpeg download error: %s", err.Error())
		}
		dl.LogInfo("FFmpeg downloaded successfully!")
	} else if err != nil {
		return fmt.Errorf("ffmpeg exist check error: %s", err.Error())
	}

	return nil
}
