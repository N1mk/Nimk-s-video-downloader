//go:build linux

package deps_downloader

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"nvd/internal/models"
	"os"
	"path/filepath"

	"github.com/ulikunitz/xz"
)

const (
	ytDlpDownloadLinkLinux  string = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux"
	ffmpegDownloadLinkLinux string = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linux64-gpl.tar.xz"
	ytDlpDownloadName       string = "yt-dlp"
	ffmpegDownloadName      string = "ffmpeg.tar.xz"
	ffmpegAppName           string = "ffmpeg"
)

func installYtDlp(ctx context.Context, exeDir string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ytDlpDownloadLinkLinux, nil)
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ffmpegDownloadLinkLinux, nil)
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
