package models

import (
	"context"
	"errors"
)

var (
	ErrDeadlineExceeded error  = errors.New("deadline exceeded")
	JobStatusInProcess  string = "in process"
	JobStatusError      string = "error"
	JobStatusComplete   string = "complete"
)

type ExtensionRequestData struct {
	URL       string `json:"url"`
	Extension string `json:"extension"`
}

type DownloadJob struct {
	Ctx       context.Context
	Link      string
	Extension string
	Status    string
	Error     error
}
