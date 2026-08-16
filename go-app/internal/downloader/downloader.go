package downloader

import (
	"bytes"
	"context"
	"fmt"
	"nvd/internal/models"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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

func (d *Downloader) GetFileName(ctx context.Context, link string) (fileName string, err error) {
	cmd := exec.CommandContext(ctx, d.ytdlpPath, "--restrict-filenames", "--get-filename", "-o", "%(title)s.%(ext)s", link)
	setPlatformSysProcAttr(cmd)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %w", models.ErrGetFilenameCommandRunError, err)
	}

	return strings.TrimSpace(out.String()), nil
}

func (d *Downloader) Download(ctx context.Context, link string, downloadPath string) (err error) {
	cmd := exec.CommandContext(ctx, d.ytdlpPath, "--restrict-filenames", "-o", "%(title)s.%(ext)s", link)
	cmd.Dir = downloadPath
	setPlatformSysProcAttr(cmd)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %w", models.ErrDownloadCommandRunError, err)
	}

	return nil
}
