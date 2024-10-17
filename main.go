package main

import (
	"flag"
	"net/http"
	"nexum/proxy"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	configFile string
	listenAddr string
	logFile    string
)

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "Path to configuration file")
	flag.StringVar(&listenAddr, "listen", ":8080", "Address to listen on")
	flag.StringVar(&logFile, "log", "proxy.log", "Path to log file")
}

var config Config

func main() {
	flag.Parse()

	initLogger()

	if err := LoadConfig(configFile); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Infof("Starting proxy server on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, http.HandlerFunc(HandleRequest)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initLogger() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
