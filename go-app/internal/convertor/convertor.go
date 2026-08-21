package convertor

import (
	"context"
	"fmt"
	"nvd/internal/project_errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Converotr interface {
	Convert(ctx context.Context, dir string, fileName string, extension string) (newFileName string, err error)
}

type DefaultConvertor struct {
	ffmpegPath string
}

func NewDefaultConvertor() (con *DefaultConvertor, err error) {
	c := &DefaultConvertor{}

	if err := c.updatePath(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *DefaultConvertor) Convert(ctx context.Context, dir string, fileName string, extension string) (newFileName string, err error) {
	if filepath.Ext(fileName) == "."+extension {
		return fileName, nil
	}

	newFileName = strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "." + extension

	fmt.Println(newFileName)

	cmd := exec.CommandContext(ctx, c.ffmpegPath, "-threads", "2", "-i", fileName, "-preset", "veryfast", newFileName)
	cmd.Dir = dir
	setPlatformSysProcAttr(cmd)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %w", project_errors.ErrConvertCommandRunError, err)
	}

	if err := os.Remove(fmt.Sprintf("%s/%s", dir, fileName)); err != nil {
		return "", fmt.Errorf("%w: %w", project_errors.ErrDeleteCommandRunError, err)
	}

	return newFileName, nil
}

type MockConvertor struct{}

func (m *MockConvertor) Convert(_ context.Context, _ string, fileName string, extension string) (newFileName string, err error) {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "." + extension, nil
}
