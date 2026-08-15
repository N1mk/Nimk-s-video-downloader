package convertor

import (
	"context"
	"fmt"
	"nvd/internal/models"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Convertor struct {
	ffmpegPath string
}

func NewConvertor() *Convertor {
	return &Convertor{}
}

func (c *Convertor) UpdatePath() error {
	var ffmpegName string
	switch runtime.GOOS {
	case "windows":
		ffmpegName = "ffmpeg.exe"
	case "linux":
		ffmpegName = "ffmpeg"
	case "darwin":
		ffmpegName = "ffmpeg"
	default:
		return models.ErrUnknownOS
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Dir(exePath)
	c.ffmpegPath = filepath.Join(dir, ffmpegName)

	return nil
}

func (c *Convertor) Convert(ctx context.Context, dir string, fileName string, extension string) error {
	if extension == "webm" {
		return nil
	}

	newFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "." + extension

	cmd := exec.CommandContext(ctx, c.ffmpegPath, "-i", fileName, newFileName)
	cmd.Dir = dir
	setPlatformSysProcAttr(cmd)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("convert command run error: %w", err)
	}

	if err := os.Remove(fmt.Sprintf("%s/%s", dir, fileName)); err != nil {
		return fmt.Errorf("delete command run error: %w", err)
	}

	return nil
}
