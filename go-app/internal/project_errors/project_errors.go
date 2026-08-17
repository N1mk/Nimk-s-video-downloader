package project_errors

import (
	"errors"
)

var (
	ErrDeadlineExceeded           error = errors.New("deadline exceeded")
	ErrNotFound                   error = errors.New("not found")
	ErrUnknownOS                  error = errors.New("unknown os")
	ErrDownloadCommandRunError    error = errors.New("download command run error")
	ErrConvertCommandRunError     error = errors.New("convert command run error")
	ErrDeleteCommandRunError      error = errors.New("delete command run error")
	ErrGetFilenameCommandRunError error = errors.New("get filename command run error")
)
