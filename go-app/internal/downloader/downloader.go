package downloader

import (
	"context"
	"fmt"
	"nvd/internal/models"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type Downloader struct {
	ytdlpPath string
}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) UpdatePath() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Dir(exePath)
	d.ytdlpPath = filepath.Join(dir, "yt-dlp.exe")

	return nil
}

func (d *Downloader) Update(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, d.ytdlpPath, "-U")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}

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
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("download command run error: %w", err)
	}

	return nil
}
