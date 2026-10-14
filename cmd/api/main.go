package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"boliplate-go/internal/config"
	"boliplate-go/internal/db"
	"boliplate-go/internal/httpserver"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		fatal("config", err)
	}
	if err := cfg.ValidateForAPI(); err != nil {
		fatal("config", err)
	}

	gormDB, err := db.Open(cfg.DB)
	if err != nil {
		fatal("db", err)
	}

	srv := httpserver.New(cfg, gormDB)

	httpSrv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           srv.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// NotifyContext cancela o ctx no primeiro SIGINT/SIGTERM: graceful shutdown sem canal manual.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("http server listening", "port", cfg.HTTPPort)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fatal("http", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown", "err", err)
	}
}

func fatal(msg string, err error) {
	slog.Error(msg, "err", err)
	os.Exit(1)
}
