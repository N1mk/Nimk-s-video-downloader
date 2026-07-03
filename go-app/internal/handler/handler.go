package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"nvd/internal/logger"
	"nvd/internal/models"
	"nvd/internal/service"
)

type ExtensionHandler struct {
	ctx context.Context
	svc *service.DownloadService
	dl  *logger.DownloaderLogger
}

func NewExtensionHandler(ctx context.Context, svc *service.DownloadService, dl *logger.DownloaderLogger) *ExtensionHandler {
	return &ExtensionHandler{ctx: ctx, svc: svc, dl: dl}
}

func (h *ExtensionHandler) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.dl.LogError(fmt.Sprint("Error method not allowed"))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var data models.ExtensionRequestData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.dl.LogError(fmt.Sprint("Error bad json"))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	h.svc.Download(h.ctx, data.URL, data.Extension)
}
