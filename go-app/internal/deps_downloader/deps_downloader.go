package deps_downloader

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"nvd/internal/logger"
	"os"
	"path/filepath"
)

func installYtDlp(installLink string, exeDir string) error {
	resp, err := http.Get(installLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	dst, err := os.Create(filepath.Join(exeDir, "yt-dlp.exe"))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, resp.Body); err != nil {
		return err
	}

	return nil
}

func installFFmpeg(installLink string, exeDir string, dl *logger.DownloaderLogger) error {
	resp, err := http.Get(installLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	tmpPath := filepath.Join(exeDir, "ffmpeg.zip")

	dl.LogInfo("Создание временного файла")
	tmp, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer os.Remove(tmpPath)
	defer tmp.Close()

	dl.LogInfo("Начало скачивания zip архива")
	if _, err := io.Copy(tmp, resp.Body); err != nil {
		return err
	}

	dl.LogInfo("Начало чтения zip архива")
	r, err := zip.OpenReader(tmpPath)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		if filepath.Base(f.Name) == "ffmpeg.exe" {
			dl.LogInfo("Найден ffmpeg file")
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			dl.LogInfo("FFmpeg file открыт")

			dst, err := os.Create(filepath.Join(exeDir, "ffmpeg.exe"))
			if err != nil {
				return err
			}
			defer dst.Close()

			dl.LogInfo("Создан ffmpeg.exe")

			if _, err := io.Copy(dst, rc); err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("yt-dlp file not found")
}

func DownloadDeps(dl *logger.DownloaderLogger) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("exe file location determining error: %s", err.Error())
	}

	exeDir := filepath.Dir(exePath)

	if _, err := os.Stat(filepath.Join(exeDir, "yt-dlp.exe")); os.IsNotExist(err) {
		dl.LogInfo("Yt-dlp not found. Downloading...")
		if err := installYtDlp("https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe", exeDir); err != nil {
			return fmt.Errorf("yt-dlp download error: %s", err.Error())
		}
		dl.LogInfo("Yt-dlp downloaded successfully!")
	} else if err != nil {
		return fmt.Errorf("Yt-dlp exist check error: %s", err.Error())
	}

	if _, err := os.Stat(filepath.Join(exeDir, "ffmpeg.exe")); os.IsNotExist(err) {
		dl.LogInfo("FFmpeg not found. Downloading...")
		if err := installFFmpeg("https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip", exeDir, dl); err != nil {
			return fmt.Errorf("ffmpeg download error: %s", err.Error())
		}
		dl.LogInfo("FFmpeg downloaded successfully!")
	} else if err != nil {
		return fmt.Errorf("FFmpeg exist check error: %s", err.Error())
	}

	return nil
}
