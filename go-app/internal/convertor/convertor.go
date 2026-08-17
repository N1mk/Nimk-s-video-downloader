package convertor

import (
	"context"
	"fmt"
	"nvd/internal/models"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Convertor struct {
	ffmpegPath string
}

func NewConvertor() *Convertor {
	return &Convertor{}
}

func (c *Convertor) Convert(ctx context.Context, dir string, fileName string, extension string) error {
	if filepath.Ext(fileName) == extension {
		pathToFile := filepath.Join(dir, fileName)
		os.Rename(pathToFile, strings.TrimSuffix(pathToFile, filepath.Ext(pathToFile))+"(downloaded)."+extension)
		return nil
	}

	newFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "." + extension

	cmd := exec.CommandContext(ctx, c.ffmpegPath, "-threads", "2", "-i", fileName, "-preset", "veryfast", newFileName)
	cmd.Dir = dir
	setPlatformSysProcAttr(cmd)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %w", models.ErrConvertCommandRunError, err)
	}

	if err := os.Remove(fmt.Sprintf("%s/%s", dir, fileName)); err != nil {
		return fmt.Errorf("%w: %w", models.ErrDeleteCommandRunError, err)
	}

	return nil
}
