package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"nvd/internal/config"
	"nvd/internal/loader"
	"nvd/internal/logger"
	"nvd/internal/service"
	"testing"
)

func TestPostJobStatusRequest(t *testing.T) {
	tests := []struct {
		name              string
		jobs              []*service.DownloadJob
		jobID             string
		expectedStatus    int
		expectedExportJob ExportJob
	}{
		{
			name:              "Ok. Job status complete",
			jobs:              []*service.DownloadJob{&service.DownloadJob{ID: 1111, Status: service.JobStatusComplete, Error: nil}},
			jobID:             `{"id": 1111}`,
			expectedStatus:    http.StatusOK,
			expectedExportJob: ExportJob{ID: 1111, Status: service.JobStatusComplete, ErrorString: ""},
		},
		{
			name:              "Ok. Job status in process",
			jobs:              []*service.DownloadJob{&service.DownloadJob{ID: 1111, Status: service.JobStatusInProcess, Error: nil}},
			jobID:             `{"id": 1111}`,
			expectedStatus:    http.StatusOK,
			expectedExportJob: ExportJob{ID: 1111, Status: service.JobStatusInProcess, ErrorString: ""},
		},
		{
			name:              "Ok. Job status already exists",
			jobs:              []*service.DownloadJob{&service.DownloadJob{ID: 1111, Status: service.JobStatusAlreadyExists, Error: nil}},
			jobID:             `{"id": 1111}`,
			expectedStatus:    http.StatusOK,
			expectedExportJob: ExportJob{ID: 1111, Status: service.JobStatusAlreadyExists, ErrorString: ""},
		},
		{
			name:              "Ok. Job status retrying",
			jobs:              []*service.DownloadJob{&service.DownloadJob{ID: 1111, Status: service.JobStatusRetrying, Error: nil}},
			jobID:             `{"id": 1111}`,
			expectedStatus:    http.StatusOK,
			expectedExportJob: ExportJob{ID: 1111, Status: service.JobStatusRetrying, ErrorString: ""},
		},
		{
			name:              "Ok. Job status error",
			jobs:              []*service.DownloadJob{&service.DownloadJob{ID: 1111, Status: service.JobStatusError, Error: fmt.Errorf("ERROR!")}},
			jobID:             `{"id": 1111}`,
			expectedStatus:    http.StatusOK,
			expectedExportJob: ExportJob{ID: 1111, Status: service.JobStatusError, ErrorString: "ERROR!"},
		},
		{
			name:              "Not found",
			jobs:              []*service.DownloadJob{&service.DownloadJob{ID: 1111, Status: service.JobStatusError, Error: fmt.Errorf("ERROR!")}},
			jobID:             `{"id": 1110}`,
			expectedStatus:    http.StatusNotFound,
			expectedExportJob: ExportJob{ID: 1111, Status: service.JobStatusError, ErrorString: "ERROR!"},
		},
	}

	svc := &service.MockDownloadService{}

	h := NewExtensionHandler(context.TODO(), svc, &logger.DownloaderLogger{}, &config.ConfigReader{}, &loader.DefaultLoader{})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "localhost:8080/status", bytes.NewBufferString(test.jobID))
			rec := httptest.NewRecorder()

			svc.Jobs = test.jobs

			h.PostJobStatusRequest(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("Wanted status code %d, got %d", test.expectedStatus, rec.Code)
			}

			var exportJob ExportJob
			if err := json.NewDecoder(rec.Body).Decode(&exportJob); err != nil {
				t.Errorf("Failed to unmarshal response body: %s", err.Error())
			}

			if exportJob != test.expectedExportJob {
				t.Errorf("Wanted job %v, got %v", test.expectedExportJob, exportJob)
			}
		})
	}
}
