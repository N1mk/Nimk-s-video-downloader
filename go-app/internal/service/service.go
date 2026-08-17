package service

import (
	"context"
	"errors"
	"fmt"
	"nvd/internal/config"
	"nvd/internal/convertor"
	"nvd/internal/downloader"
	"nvd/internal/logger"
	"nvd/internal/models"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	JobStatusInProcess     int8 = 0
	JobStatusError         int8 = 1
	JobStatusRetrying      int8 = 2
	JobStatusComplete      int8 = 3
	JobStatusAlreadyExists int8 = 4
)

type DownloadJob struct {
	ID        int             `json:"id"`
	Ctx       context.Context `json:"ctx"`
	Link      string          `json:"job.Link"`
	Extension string          `json:"extension"`
	Quality   string          `json:"quality"`
	Status    int8            `json:"status"`
	Error     error           `json:"error"`
}

type DownloadService struct {
	ctx           context.Context
	downloadPath  string
	in            chan *DownloadJob
	Wg            sync.WaitGroup
	workersCount  int
	jobs          []*DownloadJob
	dl            *logger.DownloaderLogger
	dow           *downloader.Downloader
	con           *convertor.Convertor
	maxRetryCount int
}

func NewDownloadService(ctx context.Context, downloadPath string, dl *logger.DownloaderLogger, dow *downloader.Downloader, con *convertor.Convertor, maxRetryCount int) *DownloadService {
	return &DownloadService{ctx: ctx, downloadPath: downloadPath, in: make(chan *DownloadJob), jobs: make([]*DownloadJob, 0), dl: dl, dow: dow, con: con, maxRetryCount: maxRetryCount}
}

func (s *DownloadService) CreateWorkers(amount int) {
	for i := 0; i < amount; i++ {
		s.workersCount += 1
		s.Wg.Add(1)
		go s.DownloaderWorker(s.ctx, s.in, s.workersCount)
	}
}

func (s *DownloadService) Download(ctx context.Context, id int, link string, extension string, quality string) {
	job := &DownloadJob{ID: id, Ctx: ctx, Link: link, Extension: extension, Quality: quality, Status: JobStatusInProcess}

	s.jobs = append(s.jobs, job)

	s.in <- job
}

func (s *DownloadService) GetJobByID(id int) (job *DownloadJob, err error) {
	for _, job := range s.jobs {
		if job.ID == id {
			return job, nil
		}
	}

	return nil, models.ErrNotFound
}

func (s *DownloadService) ChangeConfiguration(config *config.Config) {
	s.downloadPath = config.DownloadPath
	s.maxRetryCount = config.MaxRetryCount
}

func (s *DownloadService) DownloaderWorker(ctx context.Context, in <-chan *DownloadJob, id int) {
	for {
		select {
		case <-ctx.Done():
			s.Wg.Done()
			return
		case job := <-in:
			s.dl.LogInfo(fmt.Sprintf("Started downloading video(%s) to %s", job.Link, s.downloadPath))

			if _, err := os.Stat(s.downloadPath); os.IsNotExist(err) {
				err = os.MkdirAll(s.downloadPath, 0755)
				if err != nil {
					s.dl.LogError(fmt.Sprintf("Worker %d error: folder creation error: %s", id, err.Error()))
					job.Status = JobStatusError
					job.Error = fmt.Errorf("folder creation error: %w", err)
					continue
				}
			} else if err != nil {
				s.dl.LogError(fmt.Sprintf("Worker %d error: folder exist check error: %s", id, err.Error()))
				job.Status = JobStatusError
				job.Error = fmt.Errorf("folder exist check error: %w", err)
				continue
			}

			var fileName string

			fileName, ok, err := s.tryDownload(job)
			if !ok {
				if errors.Is(err, models.ErrDownloadCommandRunError) || errors.Is(err, models.ErrGetFilenameCommandRunError) {
					s.dl.LogError(fmt.Sprintf("Worker %d error: download error: %s", id, err.Error()))
					s.dl.LogInfo(fmt.Sprintf("Worker %d: Retrying", id))
					job.Status = JobStatusRetrying
					isDownloaded := false
					alreadyExists := false
					for i := 0; i < s.maxRetryCount; i++ {
						fileName, ok, err = s.tryDownload(job)
						if !ok {
							if err != nil {
								if errors.Is(err, models.ErrDownloadCommandRunError) || errors.Is(err, models.ErrGetFilenameCommandRunError) {
									s.dl.LogError(fmt.Sprintf("Worker %d error (Retry %d): download error: %s", id, i+1, err.Error()))
								} else {
									s.dl.LogError(fmt.Sprintf("Worker %d unknown error (Retry %d): download error: %s", id, i+1, err.Error()))
									job.Status = JobStatusError
									job.Error = err
									break
								}
							} else {
								s.dl.LogInfo(fmt.Sprintf("Worker %d error (Retry %d): video file already exists!", id, i+1))
								job.Status = JobStatusAlreadyExists
								alreadyExists = true
								break
							}
						} else if err == nil {
							s.dl.LogInfo(fmt.Sprintf("Worker %d error (Retry %d): video downloaded! Converting...", id, i+1))
							isDownloaded = true
							break
						}
					}
					if !isDownloaded || alreadyExists {
						continue
					}
				} else if err != nil {
					s.dl.LogError(fmt.Sprintf("Worker %d error: download error: %s", id, err.Error()))
					job.Status = JobStatusError
					job.Error = err
					continue
				}
				s.dl.LogInfo(fmt.Sprintf("Worker %d error: video file already exists!", id))
				job.Status = JobStatusAlreadyExists
				continue
			}

			if err := s.con.Convert(job.Ctx, s.downloadPath, fileName, job.Extension); err != nil {
				s.dl.LogError(fmt.Sprintf("Worker %d error: convert error: %s", id, err.Error()))
				job.Status = JobStatusError
				job.Error = err
				continue
			}

			job.Status = JobStatusComplete

			s.dl.LogInfo(fmt.Sprintf("Worker %d: video(%s) succesfully downloaded to %s", id, job.Link, s.downloadPath))
		}
	}
}

func (s *DownloadService) tryDownload(job *DownloadJob) (fileName string, ok bool, err error) {
	fileName, err1 := s.dow.GetFileName(job.Ctx, job.Link)

	oldExt := filepath.Ext(fileName)
	pathToFile := strings.TrimSuffix(filepath.Join(s.downloadPath, fileName), oldExt) + "." + job.Extension

	_, err = os.Stat(pathToFile)
	if err == nil && job.Extension != oldExt {
		return fileName, false, nil
	}
	err = nil

	err2 := s.dow.Download(job.Ctx, job.Link, s.downloadPath, job.Quality)

	if err1 != nil {
		err = err2
	} else if err2 != nil {
		err = err2
	}

	if err != nil {
		return "", false, err
	}

	return fileName, true, nil
}
