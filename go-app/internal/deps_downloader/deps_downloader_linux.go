//go:build linux

package deps_downloader

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"nvd/internal/logger"
	"nvd/internal/project_errors"
	"os"
	"path/filepath"

	"github.com/ulikunitz/xz"
)

const (
	ytDlpDownloadLink  string = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux"
	ffmpegDownloadLink string = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linux64-gpl.tar.xz"
	qjsDownloadLink    string = "https://github.com/quickjs-ng/quickjs/releases/latest/download/qjs-linux-x86_64"
	ytDlpDownloadName  string = "yt-dlp"
	ffmpegDownloadName string = "ffmpeg.tar.xz"
	ffmpegAppName      string = "ffmpeg"
	qjsDownloadName    string = "qjs"
)

func DownloadDeps(ctx context.Context, dl *logger.DownloaderLogger) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("exe file location determining error: %s", err.Error())
	}

	exeDir := filepath.Dir(exePath)

	if _, err := os.Stat(filepath.Join(exeDir, ytDlpDownloadName)); os.IsNotExist(err) {
		dl.LogInfo("Yt-dlp not found. Downloading...")
		if err := installYtDlp(ctx, exeDir); err != nil {
			return fmt.Errorf("yt-dlp download error: %s", err.Error())
		}
		dl.LogInfo("Yt-dlp downloaded successfully!")
	} else if err != nil {
		return fmt.Errorf("yt-dlp exist check error: %s", err.Error())
	}

	if _, err := os.Stat(filepath.Join(exeDir, ffmpegAppName)); os.IsNotExist(err) {
		dl.LogInfo("FFmpeg not found. Downloading...")
		if err := installFFmpeg(ctx, exeDir); err != nil {
			return fmt.Errorf("ffmpeg download error: %s", err.Error())
		}
		dl.LogInfo("FFmpeg downloaded successfully!")
	} else if err != nil {
		return fmt.Errorf("ffmpeg exist check error: %s", err.Error())
	}

	if _, err := os.Stat(filepath.Join(exeDir, qjsDownloadName)); os.IsNotExist(err) {
		dl.LogInfo("QuickJS not found. Downloading...")
		if err := installQJS(ctx, exeDir); err != nil {
			return fmt.Errorf("qjs download error: %s", err.Error())
		}
		dl.LogInfo("QuickJS downloaded successfully!")
	} else if err != nil {
		return fmt.Errorf("QuickJS exist check error: %s", err.Error())
	}

	return nil
}

func installYtDlp(ctx context.Context, exeDir string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ytDlpDownloadLink, nil)
	if err != nil {
		return project_errors.ErrDeadlineExceeded
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	dst, err := os.Create(filepath.Join(exeDir, ytDlpDownloadName))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, resp.Body); err != nil {
		return err
	}

	return nil
}

func installFFmpeg(ctx context.Context, exeDir string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ffmpegDownloadLink, nil)
	if err != nil {
		return project_errors.ErrDeadlineExceeded
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	tmpPath := filepath.Join(exeDir, ffmpegDownloadName)
	tmp, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer os.Remove(tmpPath)
	defer tmp.Close()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		return err
	}

	xzr, err := xz.NewReader(tmp)
	if err != nil {
		return err
	}

	r := tar.NewReader(xzr)

	for {
		header, err := r.Next()
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("ffmpeg binary not found")
		} else if err != nil {
			return fmt.Errorf("archive reader error: %w", err)
		}
		if filepath.Base(header.Name) == ffmpegAppName {
			dst, err := os.Create(filepath.Join(exeDir, ffmpegAppName))
			if err != nil {
				return err
			}
			defer dst.Close()

			if _, err := io.Copy(dst, r); err != nil {
				return err
			}

			return nil
		}
	}
}

func installQJS(ctx context.Context, exeDir string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, qjsDownloadLink, nil)
	if err != nil {
		return project_errors.ErrDeadlineExceeded
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	dst, err := os.Create(filepath.Join(exeDir, qjsDownloadName))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, resp.Body); err != nil {
		return err
	}

	if err := os.Rename(qjsDownloadLink, "qjs"); err != nil {
		return err
	}

	return nil
}
