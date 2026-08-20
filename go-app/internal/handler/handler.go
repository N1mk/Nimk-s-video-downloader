package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"nvd/internal/config"
	"nvd/internal/loader"
	"nvd/internal/logger"
	"nvd/internal/project_errors"
	"nvd/internal/service"
	"strconv"
)

type ExtensionRequest struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Extension string `json:"extension"`
	Quality   string `json:"quality"`
}

type ExportJob struct {
	ID          int    `json:"id"`
	Status      uint8  `json:"status"`
	ErrorString string `json:"error"`
}

type ExtensionHandler struct {
	ctx context.Context
	svc *service.DownloadService
	dl  *logger.DownloaderLogger
	cr  *config.ConfigReader
	loa *loader.Loader
}

func NewExtensionHandler(ctx context.Context, svc *service.DownloadService, dl *logger.DownloaderLogger, cr *config.ConfigReader, loa *loader.Loader) *ExtensionHandler {
	return &ExtensionHandler{ctx: ctx, svc: svc, dl: dl, cr: cr, loa: loa}
}

func (h *ExtensionHandler) PostDownload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data ExtensionRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.dl.LogError("Bad input")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	h.svc.Download(h.ctx, data.ID, data.URL, data.Extension, data.Quality)
}

func (h *ExtensionHandler) PostJobStatusRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data ExtensionRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.dl.LogError("Bad input")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	job, err := h.svc.GetJobByID(data.ID)
	if errors.Is(err, project_errors.ErrNotFound) {
		h.dl.LogError("Job not found")
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	exprotJob := &ExportJob{ID: job.ID, Status: job.Status, ErrorString: errStr}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(exprotJob); err != err {
		h.dl.LogError(fmt.Sprintf("Job marshaling error: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Job marshaling error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *ExtensionHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, err := h.cr.GetConfigJSON()
	if err != nil {
		h.dl.LogError(fmt.Sprintf("Config reader error: %s", err.Error()))
		http.Error(w, "Config reading error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *ExtensionHandler) PostConfig(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.dl.LogError("Bad input")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	currentConfig, err := h.cr.GetConfig()
	if err != nil {
		h.dl.LogError(fmt.Sprintf("Getting config error: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Getting config error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var rawConfig config.RawConfig
	if err := json.Unmarshal(data, &rawConfig); err != nil {
		h.dl.LogError(fmt.Sprintf("Config unmarshalling error: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Config unmarshalling error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	maxRetryCount, err := strconv.Atoi(rawConfig.MaxRetryCount)
	if err != nil {
		h.dl.LogError("Bad input")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	newConfig := config.Config{DownloadPath: rawConfig.DownloadPath, MaxRetryCount: maxRetryCount, AddToAutostart: currentConfig.AddToAutostart}

	if err := h.cr.SetConfig(&newConfig); err != nil {
		h.dl.LogError(fmt.Sprintf("Config reader error: %s", err.Error()))
		http.Error(w, "Config writing error", http.StatusInternalServerError)
		return
	}

	h.svc.ChangeConfiguration(&newConfig)
}

func (h *ExtensionHandler) PostLogs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := h.dl.OpenLogFile(); err != nil {
		h.dl.LogError(fmt.Sprintf("Logger error: %s", err.Error()))
		http.Error(w, "Cannot open log file", http.StatusInternalServerError)
	}
}

func (h *ExtensionHandler) PostDownloaderUpdate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := h.loa.Update(h.ctx); err != nil {
		h.dl.LogError(fmt.Sprintf("Downloader update error: %s", err.Error()))
		http.Error(w, "Update error", http.StatusInternalServerError)
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}
