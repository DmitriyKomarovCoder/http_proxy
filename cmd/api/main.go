package main

import (
	"context"
	"errors"
	"github.com/DmitriyKomarovCoder/http_proxy/common/closer"
	customLogger "github.com/DmitriyKomarovCoder/http_proxy/common/logger"
	"github.com/DmitriyKomarovCoder/http_proxy/config"
	handler "github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/api/delivery/http"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/api/repository"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/api/usecase"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/proxy"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const (
	pathConfig = "config"
	nameConfig = "config"
)

func main() { // TO DO: Move to internal/app
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.InitialConfig(nameConfig, pathConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := customLogger.NewLogger(cfg.Logfile.Path)
	if err != nil {
		log.Fatalf("Error logger loading: %v", err)
	}

	conn, err := repository.NewPostgresConnect(ctx, cfg.Postgres.Host,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Name,
		cfg.Postgres.Port)
	if err != nil {
		logger.Fatalf("Error connect postgres: %v", err)
	}

	repo := repository.NewRepository(conn)
	logger.Info("Db Connect successfully")

	useCase := usecase.NewUsecase(repo)
	delivery := handler.NewHandler(useCase, *logger)
	router := handler.InitRouter(delivery)

	apiServer := &http.Server{ // TO DO: move init function
		Addr:         cfg.ApiServer.Host + cfg.ApiServer.Port,
		Handler:      router,
		ReadTimeout:  cfg.ApiServer.ReadTimeout * time.Second,
		WriteTimeout: cfg.ApiServer.WriteTimeout * time.Second,
	}

	proxyHandler := proxy.NewProxy(useCase, *logger)

	proxyServer := &http.Server{ // TO DO: move init function
		Addr:         cfg.ProxyServer.Host + cfg.ProxyServer.Port,
		Handler:      proxyHandler,
		ReadTimeout:  cfg.ProxyServer.ReadTimeout * time.Second,
		WriteTimeout: cfg.ProxyServer.WriteTimeout * time.Second,
	}

	c := &closer.Closer{}
	c.Add(apiServer.Shutdown)
	c.Add(proxyServer.Shutdown)
	c.Add(func(ctx context.Context) error {
		conn.Close()
		return nil
	})

	go func() {
		if err := proxyServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Erorr starting server: %v", err)
		}
	}()

	logger.Infof("Proxy Server start in port: %v", cfg.ProxyServer.Port)

	go func() {
		if err := apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Erorr starting server: %v", err)
		}
	}()
	logger.Infof("Api Server start in port: %v", cfg.ApiServer.Port)

	<-ctx.Done()
	logger.Info("shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout*time.Second)
	defer cancel()

	if err := c.Close(shutdownCtx); err != nil {
		logger.Fatalf("closer: %v", err)
	}

	logger.Info("Service close without error")
}
