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

func (c *Convertor) Convert(ctx context.Context, dirPath string, extension string) error {
	if extension == "webm" {
		return nil
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("directory read error: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		ext := strings.ToLower(filepath.Ext(name))

		if ext == ".mp4" || ext == ".mov" || ext == ".mp3" || ext == ".aac" || ext == ".wav" || ext == "" {
			continue
		}

		outName := strings.TrimSuffix(name, filepath.Ext(name)) + "." + extension

		cmd := exec.CommandContext(ctx, c.ffmpegPath, "-i", name, outName)
		cmd.Dir = dirPath

		setPlatformSysProcAttr(cmd)

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("convert command run error: %w", err)
		}

		if err := os.Remove(fmt.Sprintf("%s/%s", dirPath, name)); err != nil {
			return fmt.Errorf("delete command run error: %w", err)
		}
	}

	return nil
}
