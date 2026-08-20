package loader

import (
	"bytes"
	"context"
	"fmt"
	"nvd/internal/project_errors"
	"os/exec"
	"strings"
)

type Loader struct {
	ytdlpPath string
	qjsPath   string
}

func NewLoader() *Loader {
	return &Loader{}
}

func (d *Loader) Update(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, d.ytdlpPath, "-U")

	setPlatformSysProcAttr(cmd)

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return project_errors.ErrDeadlineExceeded
	} else if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	return nil
}

func (d *Loader) GetFileName(ctx context.Context, link string) (fileName string, err error) {
	cmd := exec.CommandContext(ctx, d.ytdlpPath, "--restrict-filenames", "--get-filename", "--js-runtimes", fmt.Sprintf("quickjs:%s", d.qjsPath), "-o", "%(title)s.%(ext)s", link)
	setPlatformSysProcAttr(cmd)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %w", project_errors.ErrGetFilenameCommandRunError, err)
	}

	return strings.TrimSpace(out.String()), nil
}

func (d *Loader) Download(ctx context.Context, link string, downloadPath string, quality string) (err error) {
	bitrate := "16000"
	switch quality {
	case "1080":
		bitrate = "8000"
	case "720":
		bitrate = "5000"
	case "480":
		bitrate = "2500"
	}

	cmd := exec.CommandContext(ctx, d.ytdlpPath, "--restrict-filenames", "--js-runtimes", fmt.Sprintf("quickjs:%s", d.qjsPath), "-o", "%(title)s.%(ext)s", "-f", fmt.Sprintf("bv*[height<=%s][vbr<=%s]+ba/b", quality, bitrate), link)
	cmd.Dir = downloadPath
	setPlatformSysProcAttr(cmd)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %w", project_errors.ErrDownloadCommandRunError, err)
	}

	return nil
}
