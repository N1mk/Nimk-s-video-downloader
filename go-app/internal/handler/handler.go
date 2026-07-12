package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nvd/internal/config_reader"
	"nvd/internal/logger"
	"nvd/internal/models"
	"nvd/internal/service"
)

type ExtensionHandler struct {
	ctx context.Context
	svc *service.DownloadService
	dl  *logger.DownloaderLogger
	cr  *config_reader.ConfigReader
}

func NewExtensionHandler(ctx context.Context, svc *service.DownloadService, dl *logger.DownloaderLogger, cr *config_reader.ConfigReader) *ExtensionHandler {
	return &ExtensionHandler{ctx: ctx, svc: svc, dl: dl, cr: cr}
}

func (h *ExtensionHandler) PostDownload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data models.ExtensionRequestData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.dl.LogError("Bad input")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	h.svc.Download(h.ctx, data.URL, data.Extension)
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
