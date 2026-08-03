package downloader

import (
	"context"
	"fmt"
	"nvd/internal/models"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Downloader struct {
	ytdlpPath string
}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) UpdatePath() error {
	var ytDlpName string
	switch runtime.GOOS {
	case "windows":
		ytDlpName = "yt-dlp.exe"
	case "linux":
		ytDlpName = "yt-dlp_linux"
	case "darwin":
		ytDlpName = "yt-dlp_macos"
	default:
		return models.ErrUnknownOS
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Dir(exePath)
	d.ytdlpPath = filepath.Join(dir, ytDlpName)

	return nil
}

func (d *Downloader) Update(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, d.ytdlpPath, "-U")

	setPlatformSysProcAttr(cmd)

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return models.ErrDeadlineExceeded
	} else if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	return nil
}

func (d *Downloader) Download(ctx context.Context, link string, downloadPath string) error {
	cmd := exec.CommandContext(ctx, d.ytdlpPath, link)
	cmd.Dir = downloadPath

	setPlatformSysProcAttr(cmd)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("download command run error: %w", err)
	}

	return nil
}
