package models

import (
	"context"
	"errors"
)

var (
	ErrDeadlineExceeded error  = errors.New("deadline exceeded")
	ErrNotFound         error  = errors.New("not found")
	JobStatusInProcess  string = "in process"
	JobStatusError      string = "error"
	JobStatusComplete   string = "complete"
)

type ExtensionRequestData struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Extension string `json:"extension"`
}

type DownloadJob struct {
	ID        int             `json:"id"`
	Ctx       context.Context `json:"ctx"`
	Link      string          `json:"link"`
	Extension string          `json:"extension"`
	Status    string          `json:"status"`
	Error     error           `json:"error"`
}

type Config struct {
	DownloadPath string `json:"download_path"`
}
