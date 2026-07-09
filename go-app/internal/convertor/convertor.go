package convertor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

type Convertor struct {
	ffmpegPath string
}

func NewConvertor() *Convertor {
	return &Convertor{}
}

func (c *Convertor) UpdatePath() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Dir(exePath)
	c.ffmpegPath = filepath.Join(dir, "ffmpeg.exe")

	return nil
}

func (c *Convertor) Convert(ctx context.Context, dirPath string, extension string) error {
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

		if ext == ".mp4" || ext == ".mp3" || ext == "" {
			continue
		}

		outName := strings.TrimSuffix(name, filepath.Ext(name)) + "." + extension

		cmd := exec.CommandContext(ctx, c.ffmpegPath, "-i", name, outName)
		cmd.Dir = dirPath
		cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("convert command run error: %w", err)
		}

		if err := os.Remove(fmt.Sprintf("%s/%s", dirPath, name)); err != nil {
			return fmt.Errorf("delete command run error: %w", err)
		}
	}

	return nil
}
