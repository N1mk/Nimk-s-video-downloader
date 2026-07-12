package handler

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"nvd/internal/logger"
	"nvd/internal/models"
	"nvd/internal/service"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"golang.org/x/sys/windows/registry"
)

const HostName string = "com.video.downloader"
const ExtensionID string = "pfdgdanhdelhjmbelhpcoaendcdjgkpa"
const MessageTypeDownload string = "download"
const MessageTypeCheckStatus string = "check status"

type HostManifest struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Path           string   `json:"path"`
	Type           string   `json:"type"`
	AllowedOrigins []string `json:"allowed_origins"`
}

type ExtensionHandler struct {
	ctx context.Context
	svc *service.DownloadService
	dl  *logger.DownloaderLogger
}

type Message struct {
	Link        string `json:"link"`
	Extension   string `json:"extension"`
	MessageType string `json:"message type"`
}
type Response struct {
	Reply  string `json:"reply"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

func NewExtensionHandler(ctx context.Context, svc *service.DownloadService, dl *logger.DownloaderLogger) *ExtensionHandler {
	return &ExtensionHandler{ctx: ctx, svc: svc, dl: dl}
}

func (h *ExtensionHandler) InstallNativeHost() (ok bool, err error) {
	if h.IsHostInstalled() {
		return false, nil
	}

	exePath, err := os.Executable()
	if err != nil {
		return false, err
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return false, err
	}

	manifest := HostManifest{
		Name:        HostName,
		Description: "Video downloader native messaging host",
		Path:        exePath,
		Type:        "stdio",
		AllowedOrigins: []string{
			"chrome-extension://" + ExtensionID + "/",
		},
	}

	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return false, fmt.Errorf("manifset file marshaling error: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		err := h.installWindowsHost(manifestBytes)
		if err != nil {
			return false, err
		}
		return true, nil
	case "darwin":
		err := h.installUnixHost()
		if err != nil {
			return false, err
		}
		return true, nil
	case "linux":
		err := h.installUnixHost()
		if err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, fmt.Errorf("unknown os")
	}
}

func (h *ExtensionHandler) IsHostInstalled() bool {
	switch runtime.GOOS {
	case "windows":
		keyPath := `Software\Google\Chrome\NativeMessagingHosts\` + HostName
		key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
		if err != nil {
			if err == registry.ErrNotExist {
				return false
			}
			h.dl.LogError(fmt.Sprintf("Host installation check error: %s", err.Error()))
			return false
		}
		defer key.Close()

		return true
	case "darwin":
		h.dl.LogError("No realization of host installer for mac os in this version!")
		return false
	case "linux":
		h.dl.LogError("No realization of host installer for linux in this version!")
		return false
	default:
		h.dl.LogError("Unknown OS!")
		return false
	}
}

func (h *ExtensionHandler) installWindowsHost(manifestBytes []byte) error {
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	manifestPath := filepath.Join(exeDir, HostName+".json")

	err := os.WriteFile(manifestPath, manifestBytes, 0644)
	if err != nil {
		return fmt.Errorf("manifest file creation error: %w", err)
	}

	keyPath := `Software\Google\Chrome\NativeMessagingHosts\` + HostName
	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("windows registry key creation err: %w", err)
	}
	defer key.Close()

	return key.SetStringValue("", manifestPath)
}

func (h *ExtensionHandler) installUnixHost() error {
	return fmt.Errorf("no realization of host installer for unix in this version")
}

func (h *ExtensionHandler) Start() error {
	for {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-h.ctx.Done():
			return nil
		default:
			var length uint32
			err := binary.Read(os.Stdin, binary.LittleEndian, &length)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			msgBytes := make([]byte, length)
			_, err = io.ReadFull(os.Stdin, msgBytes)
			if err != nil {
				h.dl.LogError(fmt.Sprintf("message reading error: %s", err.Error()))
			}

			var msg Message
			json.Unmarshal(msgBytes, &msg)

			if msg.MessageType == MessageTypeDownload {
				h.svc.Download(h.ctx, msg.Link, msg.Extension)

				resp := Response{Reply: "Video downloading...", Status: models.JobStatusInProcess, Error: "nil"}
				respBytes, err := json.Marshal(resp)
				if err != nil {
					h.dl.LogError(fmt.Sprintf("response marshaling error: %s", err.Error()))
				}

				binary.Write(os.Stdout, binary.LittleEndian, uint32(len(respBytes)))

				os.Stdout.Write(respBytes)
			} else if msg.MessageType == MessageTypeCheckStatus {
				var resp Response

				job, err := h.svc.GetJobByLink(msg.Link)
				if err != nil {
					resp = Response{Reply: "No jobs with this link", Status: models.JobStatusError, Error: err.Error()}
				} else {
					resp = Response{Reply: "Downloading status:", Status: job.Status, Error: job.Error.Error()}
				}

				respBytes, err := json.Marshal(resp)
				if err != nil {
					h.dl.LogError(fmt.Sprintf("response marshaling error: %s", err.Error()))
				}

				binary.Write(os.Stdout, binary.LittleEndian, uint32(len(respBytes)))

				os.Stdout.Write(respBytes)
			}
		}
	}
}
