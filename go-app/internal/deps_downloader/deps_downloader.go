package deps_downloader

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"nvd/internal/logger"
	"nvd/internal/models"
	"os"
	"path/filepath"
	"runtime"
)

const (
	YtDlpDownloadLinkWin    string = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
	FFmpegDownloadLinkWin   string = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip"
	YtDlpDownloadLinkLinux  string = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux"
	FFmpegDownloadLinkLinux string = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linux64-gpl.tar.xz"
	YtDlpDownloadLinkMacOS  string = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos"
	FFmpegDownloadLinkMacOS string = "https://evermeet.cx/ffmpeg/ffmpeg-8.1.2.zip"
)

func installYtDlp(ctx context.Context, installLink string, exeDir string) error {
	var downloadName string

	switch runtime.GOOS {
	case "windows":
		downloadName = "yt-dlp.exe"
	case "linux":
		downloadName = "yt-dlp_linux"
	case "darwin":
		downloadName = "yt-dlp_macos"
	default:
		return fmt.Errorf("unkonown OS")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, installLink, nil)
	if err != nil {
		return models.ErrDeadlineExceeded
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	dst, err := os.Create(filepath.Join(exeDir, downloadName))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, resp.Body); err != nil {
		return err
	}

	return nil
}

func installFFmpeg(ctx context.Context, installLink string, exeDir string) error {
	var downloadName string
	var appName string
	switch runtime.GOOS {
	case "windows":
		downloadName = "ffmpeg.zip"
		appName = "ffmpeg.exe"
	case "linux":
		downloadName = "ffmpeg.tar.xz"
		appName = "ffmpeg"
	case "darwin":
		downloadName = "ffmpeg.zip"
		appName = "ffmpeg"
	default:
		return fmt.Errorf("unkonown OS")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, installLink, nil)
	if err != nil {
		return models.ErrDeadlineExceeded
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	tmpPath := filepath.Join(exeDir, downloadName)
	tmp, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer os.Remove(tmpPath)
	defer tmp.Close()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		return err
	}

	r, err := zip.OpenReader(tmpPath)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		if filepath.Base(f.Name) == appName {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			dst, err := os.Create(filepath.Join(exeDir, appName))
			if err != nil {
				return err
			}
			defer dst.Close()

			if _, err := io.Copy(dst, rc); err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("yt-dlp file not found")
}

func DownloadDeps(ctx context.Context, dl *logger.DownloaderLogger) error {
	var ytDlpDownloadLink string
	var ffmpegDownloadLink string

	switch runtime.GOOS {
	case "windows":
		ytDlpDownloadLink = YtDlpDownloadLinkWin
		ffmpegDownloadLink = FFmpegDownloadLinkWin
	case "linux":
		ytDlpDownloadLink = YtDlpDownloadLinkLinux
		ffmpegDownloadLink = FFmpegDownloadLinkLinux
	case "darwin":
		ytDlpDownloadLink = YtDlpDownloadLinkMacOS
		ffmpegDownloadLink = FFmpegDownloadLinkMacOS
	default:
		return fmt.Errorf("unkonown OS")
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("exe file location determining error: %s", err.Error())
	}

	exeDir := filepath.Dir(exePath)

	if _, err := os.Stat(filepath.Join(exeDir, "yt-dlp.exe")); os.IsNotExist(err) {
		dl.LogInfo("Yt-dlp not found. Downloading...")
		if err := installYtDlp(ctx, ytDlpDownloadLink, exeDir); err != nil {
			return fmt.Errorf("yt-dlp download error: %s", err.Error())
		}
		dl.LogInfo("Yt-dlp downloaded successfully!")
	} else if err != nil {
		return fmt.Errorf("yt-dlp exist check error: %s", err.Error())
	}

	if _, err := os.Stat(filepath.Join(exeDir, "ffmpeg.exe")); os.IsNotExist(err) {
		dl.LogInfo("FFmpeg not found. Downloading...")
		if err := installFFmpeg(ctx, ffmpegDownloadLink, exeDir); err != nil {
			return fmt.Errorf("ffmpeg download error: %s", err.Error())
		}
		dl.LogInfo("FFmpeg downloaded successfully!")
	} else if err != nil {
		return fmt.Errorf("ffmpeg exist check error: %s", err.Error())
	}

	return nil
}
