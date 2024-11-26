package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"main.go/internal/config"
	"main.go/internal/handler"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting application")

	shutdownCh := make(chan struct{})
	port := ":" + fmt.Sprintf("%d", cfg.Port)
	handler.ListenPort(port, cfg.TlsPath, shutdownCh, log, cfg.Timeout)
	handler.ListenStopSig()
	close(shutdownCh)
	time.Sleep(cfg.Timeout)

}

// setup level of logger info
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
