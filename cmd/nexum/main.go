package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"nexum/internal/config"
	"nexum/internal/logger"
	"nexum/internal/proxy"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	listenAddr := flag.String("listen", ":8080", "Address to listen on")
	logFile := flag.String("log", "proxy.log", "Path to log file")
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger, err := logger.New(*logFile)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	server := proxy.NewServer(cfg, logger)

	httpServer := &http.Server{
		Addr:    *listenAddr,
		Handler: server,
	}

	go func() {
		logger.Info("Starting proxy server on %s", *listenAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: %v", err)
		}
	}()

	// gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exiting")
}
