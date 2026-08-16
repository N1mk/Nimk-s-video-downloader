//go:build windows

package deps_downloader

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"nvd/internal/models"
	"os"
	"path/filepath"
)

const (
	ytDlpDownloadLink  string = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
	ffmpegDownloadLink string = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip"
	ytDlpDownloadName  string = "yt-dlp.exe"
	ffmpegDownloadName string = "ffmpeg.zip"
	ffmpegAppName      string = "ffmpeg.exe"
)

func installYtDlp(ctx context.Context, exeDir string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ytDlpDownloadLink, nil)
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

	r, err := zip.OpenReader(tmpPath)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		if filepath.Base(f.Name) == ffmpegAppName {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			dst, err := os.Create(filepath.Join(exeDir, ffmpegAppName))
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
