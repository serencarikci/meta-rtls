package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/denizyetis/meta-rtls/internal/app"
	"github.com/denizyetis/meta-rtls/internal/config"
	"github.com/denizyetis/meta-rtls/internal/platform/db"
	"github.com/denizyetis/meta-rtls/internal/platform/logging"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "err", err)
		os.Exit(1)
	}

	logger, closeLogs, err := logging.New(cfg.AppEnv)
	if err != nil {
		slog.Error("logger setup failed", "err", err)
		os.Exit(1)
	}
	defer closeLogs()
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	oracleDB, err := db.OpenOracle(cfg.OracleDSN())
	if err != nil {
		logger.Error("oracle connection failed", "err", err)
		os.Exit(1)
	}
	defer oracleDB.Close()

	application, err := app.New(cfg, oracleDB, logger)
	if err != nil {
		logger.Error("app bootstrap failed", "err", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           application.Router(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		logger.Info("MetaRTLS API listening", "port", cfg.AppPort, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	application.Close()
}
