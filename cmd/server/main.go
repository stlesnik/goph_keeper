package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/stlesnik/goph_keeper/internal/config"
	"github.com/stlesnik/goph_keeper/internal/logger"
	"github.com/stlesnik/goph_keeper/internal/server"
	"github.com/stlesnik/goph_keeper/internal/store"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// init config
	cfg, err := config.LoadServerConfig()
	if err != nil {
		panic(err)
	}
	if cfg.PostgresDSN == "" {
		panic(fmt.Errorf("postgres DSN not set"))
	}

	err = logger.InitLogger(cfg.Environment)
	if err != nil {
		panic(err)
	}

	st, err := store.NewStore(cfg)
	if err != nil {
		logger.Logger.Fatal(err)
	}

	srv, err := server.NewServer(&cfg, st)
	if err != nil {
		logger.Logger.Fatal(err)
	}

	httpErrCh := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			httpErrCh <- err
		} else {
			close(httpErrCh)
		}
	}()
	logger.Logger.Infof("HTTP сервер запущен на %s", cfg.ServerAddress)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case sig := <-sigCh:
		logger.Logger.Infow("Получен сигнал завершения", "signal", sig)
	case err := <-httpErrCh:
		if err != nil {
			logger.Logger.Fatal("Ошибка HTTP сервера: ", "error", err)
		}
		return
	}

	logger.Logger.Info("Завершение работы серверов...")

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	if err := srv.Stop(shutdownCtx); err != nil {
		logger.Logger.Fatal("Ошибка при остановке HTTP сервера", "error", err)
	}

	if err := st.Close(); err != nil {
		logger.Logger.Fatal("Ошибка при разрыве связи с бд", "error", err)
	}
}
