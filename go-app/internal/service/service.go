package service

import (
	"context"
	"fmt"
	"nvd/internal/convertor"
	"nvd/internal/downloader"
	"nvd/internal/logger"
	"nvd/internal/models"
	"os"
	"sync"
)

type DownloadService struct {
	ctx          context.Context
	downloadPath string
	in           chan models.DownloadJob
	Wg           sync.WaitGroup
	workersCount int
	dl           *logger.DownloaderLogger
	dow          *downloader.Downloader
	con          *convertor.Convertor
}

func NewDownloadService(ctx context.Context, downloadPath string, dl *logger.DownloaderLogger, dow *downloader.Downloader, con *convertor.Convertor) *DownloadService {
	return &DownloadService{ctx: ctx, downloadPath: downloadPath, in: make(chan models.DownloadJob), dl: dl, dow: dow, con: con}
}

func (s *DownloadService) CreateWorkers(amount int) {
	for i := 0; i < amount; i++ {
		s.workersCount += 1
		s.Wg.Add(1)
		go s.DownloaderWorker(s.ctx, s.in, s.workersCount)
	}
}

func (s *DownloadService) Download(ctx context.Context, link string, extension string) {
	job := models.DownloadJob{Ctx: ctx, Link: link, Extension: extension, Status: models.JobStatusInProcess}

	s.in <- job
}

func (s *DownloadService) ChangeDownloadPath(path string) {
	s.downloadPath = path
}

func (s *DownloadService) DownloaderWorker(ctx context.Context, in <-chan models.DownloadJob, id int) {
	for {
		select {
		case <-ctx.Done():
			s.Wg.Done()
			return
		case job := <-in:
			jobCtx := job.Ctx
			link := job.Link
			extension := job.Extension

			s.dl.LogInfo(fmt.Sprintf("Started downloading video(%s) to %s", link, s.downloadPath))

			if _, err := os.Stat(s.downloadPath); os.IsNotExist(err) {
				err = os.MkdirAll(s.downloadPath, 0755)
				if err != nil {
					s.dl.LogError(fmt.Sprintf("worker %d error: folder creation error: %s", id, err.Error()))
					job.Status = models.JobStatusError
					job.Error = fmt.Errorf("folder creation error: %w", err)
				}
			} else if err != nil {
				s.dl.LogError(fmt.Sprintf("worker %d error: folder exist check error: %s", id, err.Error()))
				job.Status = models.JobStatusError
				job.Error = fmt.Errorf("folder exist check error: %w", err)
			}

			if err := s.dow.Download(jobCtx, link, s.downloadPath); err != nil {
				s.dl.LogError(fmt.Sprintf("worker %d error: downloader error: %s", id, err.Error()))
				job.Status = models.JobStatusError
				job.Error = err
			}

			if err := s.con.Convert(jobCtx, s.downloadPath, extension); err != nil {
				s.dl.LogError(fmt.Sprintf("worker %d error: convertor error: %s", id, err.Error()))
				job.Status = models.JobStatusError
				job.Error = err
			}

			job.Status = models.JobStatusComplete

			s.dl.LogInfo(fmt.Sprintf("Video(%s) succesfully downloaded to %s", link, s.downloadPath))
		}
	}
}
