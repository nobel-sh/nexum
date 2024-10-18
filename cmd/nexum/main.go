package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"nexum/internal/config"
	"nexum/internal/logger"
	"nexum/internal/proxy"
	"nexum/internal/rules"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	rulesFile := flag.String("rules", "rules.yaml", "Path to configuration file")
	flag.Parse()

	serverCfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	rulesCfg, err := rules.LoadRules(*rulesFile)
	if err != nil {
		log.Fatalf("Failed to load rules: %v", err)
	}

	logger, err := logger.New(serverCfg.LogFile)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	server := proxy.NewServer(rulesCfg, logger)

	httpServer := &http.Server{
		Addr:    serverCfg.ListenAddr,
		Handler: server,
	}

	go func() {
		logger.Info("Starting proxy server on %s", serverCfg.ListenAddr)
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
