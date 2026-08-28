package main

import (
	"context"
	"fmt"
	"syscall"

	"net/http"
	"nvd/internal/autostarter"
	"nvd/internal/config"
	"nvd/internal/convertor"
	"nvd/internal/deps_downloader"
	"nvd/internal/handler"
	"nvd/internal/loader"
	"nvd/internal/logger"
	"os/signal"
	"path/filepath"
	"time"

	"nvd/internal/service"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("FATAL: Work directory check error: %s\n", err.Error())
		os.Exit(1)
	}
	os.Chdir(filepath.Dir(exePath))

	dl, err := logger.InitDownloaderLogger()
	if err != nil {
		fmt.Printf("Logger initialization error: %s\n", err.Error())
		os.Exit(1)
	}

	cr, err := config.NewConfigReader(exePath)
	if err != nil {
		dl.LogFatal(fmt.Sprintf("Config reader error: %s", err.Error()))
		os.Exit(1)
	}

	config, err := cr.GetConfig()
	if err != nil {
		dl.LogFatal(fmt.Sprintf("Config reader error: %s", err.Error()))
		os.Exit(1)
	}

	depsDowCtx, closeDepsCtx := context.WithTimeout(context.Background(), 10*time.Minute)
	defer closeDepsCtx()

	if err := deps_downloader.DownloadDeps(depsDowCtx, dl); err != nil {
		dl.LogFatal(fmt.Sprintf("Dependencies download error: %s", err.Error()))
		os.Exit(1)
	}
	dl.LogInfo("All dependencies are installed!")

	if config.AddToAutostart {
		if ok, err := autostarter.AddToAutostart(dl); !ok {
			if err != nil {
				dl.LogFatal(fmt.Sprintf("Add to autostart error: %s", err.Error()))
				os.Exit(1)
			}
			dl.LogInfo("Program was in autostart")
		} else {
			dl.LogInfo("Program added to autostart")
		}
	}

	loa, err := loader.NewDefaultLoader()
	if err != nil {
		dl.LogFatal(fmt.Sprintf("Loader initialization error: %s", err.Error()))
		os.Exit(1)
	}

	con, err := convertor.NewDefaultConvertor()
	if err != nil {
		dl.LogFatal(fmt.Sprintf("Convertor initialization error: %s", err.Error()))
		os.Exit(1)
	}

	ctx, close := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer close()

	svc := service.NewDefaultDownloadService(ctx, config.DownloadPath, dl, loa, con, config.MaxRetryCount)

	svc.CreateWorkers(5)

	h := handler.NewExtensionHandler(ctx, svc, dl, cr, loa)

	r := chi.NewRouter()

	r.Use(handler.CorsMiddleware)
	r.Use(dl.Middleware)

	r.Post("/download", h.PostDownload)
	r.Post("/status", h.PostJobStatusRequest)
	r.Get("/config", h.GetConfig)
	r.Post("/config", h.PostConfig)
	r.Post("/logs", h.PostLogs)
	r.Post("/update", h.PostDownloaderUpdate)

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.WriteHeader(http.StatusOK)
	})

	go http.ListenAndServe("localhost:8080", r)

	svc.Wg.Wait()
}
