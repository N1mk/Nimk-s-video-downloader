package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"nvd/internal/config_reader"
	"nvd/internal/downloader"
	"nvd/internal/logger"
	"nvd/internal/models"
	"nvd/internal/service"
)

type ExtensionHandler struct {
	ctx context.Context
	svc *service.DownloadService
	dl  *logger.DownloaderLogger
	cr  *config_reader.ConfigReader
	dow *downloader.Downloader
}

func NewExtensionHandler(ctx context.Context, svc *service.DownloadService, dl *logger.DownloaderLogger, cr *config_reader.ConfigReader, dow *downloader.Downloader) *ExtensionHandler {
	return &ExtensionHandler{ctx: ctx, svc: svc, dl: dl, cr: cr, dow: dow}
}

func (h *ExtensionHandler) PostDownload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data models.ExtensionRequestData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.dl.LogError("Bad input")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	h.svc.Download(h.ctx, data.ID, data.URL, data.Extension)
}

func (h *ExtensionHandler) PostJobStatusRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data models.ExtensionRequestData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.dl.LogError("Bad input")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	job, err := h.svc.GetJobByID(data.ID)
	if errors.Is(err, models.ErrNotFound) {
		h.dl.LogError("Job not found")
		http.Error(w, "Job not found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(job); err != err {
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
	}

	if err := h.cr.SetConfigJSON(data); err != nil {
		h.dl.LogError(fmt.Sprintf("Config reader error: %s", err.Error()))
		http.Error(w, "Config writing error", http.StatusInternalServerError)
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		h.dl.LogError("Bad input")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	h.svc.ChangeDownloadPath(config.DownloadPath)
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

	if err := h.dow.Update(h.ctx); err != nil {
		h.dl.LogError(fmt.Sprintf("Downloader update error: %s", err.Error()))
		http.Error(w, "Update error", http.StatusInternalServerError)
	}
}

func (h *ExtensionHandler) Options(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}
