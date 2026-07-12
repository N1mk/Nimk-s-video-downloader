package main

import (
	"context"
	"fmt"

	//"net/http"
	"nvd/internal/autostarter"
	"nvd/internal/convertor"
	"nvd/internal/deps_downloader"
	"nvd/internal/downloader"
	"nvd/internal/handler"
	"nvd/internal/logger"
	"os/signal"
	"path/filepath"
	"time"

	"nvd/internal/service"
	"os"

	//"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Work directory check error: %s", err.Error())
	}
	os.Chdir(filepath.Dir(exePath))

	dl, err := logger.InitDownloaderLogger("app.log")
	if err != nil {
		fmt.Printf("Logger initialization error: %s", err.Error())
		os.Exit(1)
	}

	depsDowCtx, closeDepsCtx := context.WithTimeout(context.Background(), 10*time.Minute)
	defer closeDepsCtx()

	if err := deps_downloader.DownloadDeps(depsDowCtx, dl); err != nil {
		dl.LogError(fmt.Sprintf("Dependencies download error: %s", err.Error()))
	}
	dl.LogInfo("All dependencies are installed!")

	if ok, err := autostarter.AddToAutostart(); !ok {
		if err != nil {
			dl.LogFatal(fmt.Sprintf("Add to autostart error: %s", err.Error()))
			os.Exit(1)
		}
		dl.LogInfo("Program was in autostart")
	} else {
		dl.LogInfo("Program added to autostart")
	}

	dow := downloader.NewDownloader()

	if err := dow.UpdatePath(); err != nil {
		dl.LogFatal(fmt.Sprintf("Downloader path update error: %s", err.Error()))
		os.Exit(1)
	}

	dowUpdCtx, closeDowCtx := context.WithTimeout(context.Background(), 10*time.Minute)
	defer closeDowCtx()

	if err := dow.Update(dowUpdCtx); err != nil {
		dl.LogFatal(fmt.Sprintf("Downloader update error: %s", err.Error()))
		os.Exit(1)
	}

	con := convertor.NewConvertor()
	if err := con.UpdatePath(); err != nil {
		dl.LogFatal(fmt.Sprintf("Convertor path update error: %s", err.Error()))
		os.Exit(1)
	}

	if err := godotenv.Load("./config.env"); err != nil {
		dl.LogFatal(fmt.Sprintf("Env file reed error: %s", err.Error()))
		os.Exit(1)
	}

	downloadPath := os.Getenv("DOWNLOAD_PATH")

	ctx, close := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer close()

	svc := service.NewDownloadService(ctx, downloadPath, dl, dow, con)

	svc.CreateWorkers(5)

	h := handler.NewExtensionHandler(ctx, svc, dl)

	if ok, err := h.InstallNativeHost(); !ok {
		if err != nil {
			dl.LogError(fmt.Sprintf("Native host installation error: %s", err.Error()))
		}
		dl.LogInfo("Native host was created")
	} else {
		dl.LogInfo("Native host created")
	}

	if err := h.Start(); err != nil {
		dl.LogFatal(fmt.Sprintf("Handler error: %s", err.Error()))
		if ctx.Err() == nil {
			close()
		}
	}

	//<-ctx.Done()
	svc.Wg.Wait()
}
