package main

import (
	"context"
	"fmt"
	"net/http"
	"nvd/internal/autostarter"
	"nvd/internal/convertor"
	"nvd/internal/deps_downloader"
	"nvd/internal/downloader"
	"nvd/internal/handler"
	"nvd/internal/logger"
	"path/filepath"

	"nvd/internal/service"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type RequestData struct {
	URL string `json:"url"`
}

func main() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("work directory change error: %s", err.Error())
	}
	os.Chdir(filepath.Dir(exePath))

	dl, err := logger.InitDownloaderLogger("app.log")
	if err != nil {
		fmt.Printf("logger initialization error: %s", err.Error())
		os.Exit(1)
	}

	if err := deps_downloader.DownloadDeps(dl); err != nil {
		dl.LogError(fmt.Sprintf("dependencies download error: %s", err.Error()))
	}
	dl.LogInfo("All dependencies are installed!")

	if ok, err := autostarter.AddToAutostart(); !ok {
		if err != nil {
			dl.LogError(fmt.Sprintf("add to autostart error: %s", err.Error()))
			os.Exit(1)
		}
		dl.LogInfo("Program is in autostart")
	} else {
		dl.LogInfo("Program added to autostart")
	}

	dow := downloader.NewDownloader()

	if err := dow.UpdatePath(context.TODO()); err != nil {
		dl.LogError(fmt.Sprintf("downloader path update error: %s", err.Error()))
		os.Exit(1)
	}

	if err := dow.Update(context.TODO()); err != nil {
		dl.LogError(fmt.Sprintf("downloader update error: %s", err.Error()))
		os.Exit(1)
	}

	con := convertor.NewConvertor()
	if err := con.UpdatePath(context.Background()); err != nil {
		dl.LogError(fmt.Sprintf("convertor path update error: %s", err.Error()))
		os.Exit(1)
	}

	if err := godotenv.Load("./config.env"); err != nil {
		dl.LogError(fmt.Sprintf("env file reed error: %s", err.Error()))
		os.Exit(1)
	}

	downloadPath := os.Getenv("DOWNLOAD_PATH")

	svc := service.NewDownloadService(downloadPath, dl, dow, con)

	svc.CreateWorkers(5)

	h := handler.NewExtensionHandler(context.TODO(), svc, dl)

	r := chi.NewRouter()

	r.Post("/", h.Post)

	http.ListenAndServe("localhost:8080", r)
}
