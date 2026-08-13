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

func (c *Convertor) Convert(ctx context.Context, dirPath string, neededExt string) error {
	if neededExt == "webm" {
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

		currentName := file.Name()
		currentExt := strings.ToLower(filepath.Ext(currentName))

		if currentExt == ".mp4" || currentExt == ".mov" || currentExt == ".mp3" || currentExt == ".aac" || currentExt == ".wav" || currentExt == "" {
			continue
		}

		outName := strings.TrimSuffix(currentName, filepath.Ext(currentName)) + "." + neededExt

		cmd := exec.CommandContext(ctx, c.ffmpegPath, "-i", currentName, outName)
		cmd.Dir = dirPath

		setPlatformSysProcAttr(cmd)

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("convert command run error: %w", err)
		}

		if err := os.Remove(fmt.Sprintf("%s/%s", dirPath, currentName)); err != nil {
			return fmt.Errorf("delete command run error: %w", err)
		}
	}

	return nil
}
