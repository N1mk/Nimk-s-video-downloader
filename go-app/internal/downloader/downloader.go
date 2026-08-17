package downloader

import (
	"bytes"
	"context"
	"fmt"
	"nvd/internal/models"
	"os/exec"
	"strings"
)

type Downloader struct {
	ytdlpPath string
}

func NewDownloader() *Downloader {
	return &Downloader{}
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

func (d *Downloader) Download(ctx context.Context, link string, downloadPath string, quality string) (err error) {
	fmt.Println(quality) // delete

	cmd := exec.CommandContext(ctx, d.ytdlpPath, "--restrict-filenames", "-o", "%(title)s.%(ext)s", "-f", fmt.Sprintf("bv*[height<=%s]+ba/b", quality), link)
	cmd.Dir = downloadPath
	setPlatformSysProcAttr(cmd)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %w", models.ErrDownloadCommandRunError, err)
	}

	return nil
}
