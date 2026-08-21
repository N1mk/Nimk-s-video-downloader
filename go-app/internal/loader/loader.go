package loader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"nvd/internal/project_errors"
	"os/exec"
)

type Loader interface {
	Update(ctx context.Context) error
	GetFileName(ctx context.Context, link string) (fileName string, err error)
	Download(ctx context.Context, link string, downloadPath string, quality string) error
}

type DefaultLoader struct {
	ytdlpPath string
	qjsPath   string
}

func NewDefaultLoader() (loa *DefaultLoader, err error) {
	l := &DefaultLoader{}

	if err := l.updatePath(); err != nil {
		return nil, err
	}

	return l, nil
}

func (l *DefaultLoader) Update(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, l.ytdlpPath, "-U")

	setPlatformSysProcAttr(cmd)

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return project_errors.ErrDeadlineExceeded
	} else if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	return nil
}

func (l *DefaultLoader) GetFileName(ctx context.Context, link string) (fileName string, err error) {
	cmd := exec.CommandContext(ctx, l.ytdlpPath, "--print", "%()j", "--js-runtimes", fmt.Sprintf("quickjs:%s", l.qjsPath), link)
	setPlatformSysProcAttr(cmd)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %w", project_errors.ErrGetFilenameCommandRunError, err)
	}

	var meta struct {
		Title string `json:"title"`
		Ext   string `json:"ext"`
	}

	if err := json.Unmarshal(out.Bytes(), &meta); err != nil {
		return "", fmt.Errorf("json parse error: %w", err)
	}

	return meta.Title + "." + meta.Ext, nil
}

func (l *DefaultLoader) Download(ctx context.Context, link string, downloadPath string, quality string) error {
	bitrate := "35000"
	switch quality {
	case "1440":
		bitrate = "16000"
	case "1080":
		bitrate = "8000"
	case "720":
		bitrate = "5000"
	case "480":
		bitrate = "2500"
	case "360":
		bitrate = "1500"
	}

	cmd := exec.CommandContext(ctx, l.ytdlpPath, "--js-runtimes", fmt.Sprintf("quickjs:%s", l.qjsPath), "--no-restrict-filenames", "-o", "%(title)s.%(ext)s", "-f", fmt.Sprintf("bv*[height<=%s][vbr<=%s]+ba/b", quality, bitrate), link)
	cmd.Dir = downloadPath
	setPlatformSysProcAttr(cmd)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %w", project_errors.ErrDownloadCommandRunError, err)
	}

	return nil
}

type MockLoader struct{}

func (m *MockLoader) Update(_ context.Context) error {
	return nil
}

func (m *MockLoader) GetFileName(_ context.Context, _ string) (fileName string, err error) {
	return "downloaded_video.webm", nil
}

func (m *MockLoader) Download(_ context.Context, _ string, _ string, _ string) error {
	return nil
}
