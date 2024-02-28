package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/DmitriyKomarovCoder/http_proxy/common/closer"
	customLogger "github.com/DmitriyKomarovCoder/http_proxy/common/logger"
	"github.com/DmitriyKomarovCoder/http_proxy/config"
	handler "github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/api/delivery/http"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/api/repository"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/api/usecase"
	"github.com/DmitriyKomarovCoder/http_proxy/internal/pkg/customproxy"
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
	fmt.Println(cfg.Logfile.Path)
	logger, err := customLogger.NewLogger(cfg.Logfile.Path)
	if err != nil {
		log.Fatalf("Error logger loading: %v", err)
	}

	cReq, cRes, client, err := repository.NewMongoConnect(ctx, cfg.MongoDB.Url, cfg.MongoDB.DbName, cfg.MongoDB.ColRequest, cfg.MongoDB.ColResponse)
	if err != nil {
		logger.Fatalf("Error connect mongoDb: %v", err)
	}

	repo := repository.NewRepository(cReq, cRes)
	logger.Info("Db Connect successfully")

	useCase := usecase.NewUsecase(repo)
	delivery := handler.NewHandler(useCase, *logger)
	router := handler.InitRouter(delivery)

	apiServer := &http.Server{
		Addr:         cfg.ApiServer.Host + ":" + cfg.ApiServer.Port,
		Handler:      router,
		ReadTimeout:  cfg.ApiServer.ReadTimeout * time.Second,
		WriteTimeout: cfg.ApiServer.WriteTimeout * time.Second,
	}

	ca, err := customproxy.LoadCA(cfg.Certificate.Cert, cfg.Certificate.Key, cfg.Certificate.Subject)
	if err != nil {
		logger.Fatalf("error download CA: %v", err)
	}

	proxyHandler := customproxy.NewProxy(useCase, *logger, &ca)

	proxyServer := &http.Server{
		Addr:         cfg.ProxyServer.Host + ":" + cfg.ProxyServer.Port,
		Handler:      proxyHandler,
		ReadTimeout:  cfg.ProxyServer.ReadTimeout * time.Second,
		WriteTimeout: cfg.ProxyServer.WriteTimeout * time.Second,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)), // Disable HTTP/2.
	}

	c := &closer.Closer{}
	c.Add(apiServer.Shutdown)
	c.Add(proxyServer.Shutdown)
	c.Add(func(ctx context.Context) error {
		err := client.Disconnect(ctx)
		return err
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
