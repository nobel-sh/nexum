package main

import (
	"flag"
	"log"

	"nexum/internal/config"
	"nexum/internal/logger"
	"nexum/internal/proxy"
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
	logger.Info("Starting proxy server on %s", *listenAddr)
	if err := server.ListenAndServe(*listenAddr); err != nil {
		logger.Fatal("Failed to start server: %v", err)
	}
}
