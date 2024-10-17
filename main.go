// main.go
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"nexum/config"
	"nexum/rules"
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

func main() {
	flag.Parse()

	initLogger()

	if err := config.LoadConfig(configFile); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Infof("Starting proxy server on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, http.HandlerFunc(handleRequest)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	logEntry := fmt.Sprintf("[%s] %s %s %s\n",
		startTime.Format(time.RFC3339),
		r.RemoteAddr,
		r.Method,
		r.URL.String())

	log.Infof("Received request %s %s", r.Method, r.URL.String())

	rule := rules.MatchRule(config.GetRules(), r.URL.String())
	if rule != nil {
		switch rule.Action {
		case "block":
			http.Error(w, "Access denied", http.StatusForbidden)
			log.Infof("Blocked request %s based on rule", r.URL.String())
			logEntry += fmt.Sprintf("    Result: Blocked (Rule: %s)\n", rule.URLPattern)
			writeLogEntry(logEntry)
			return
		case "modify":
			rules.ApplyModifications(r, rule.Modifications)
			log.Infof("Modified request %s based on rule", r.URL.String())
			logEntry += fmt.Sprintf("    Result: Modified (Rule: %s)\n", rule.URLPattern)
		case "allow":
			log.Infof("Allowed request %s based on rule", r.URL.String())
			logEntry += fmt.Sprintf("    Result: Allowed (Rule: %s)\n", rule.URLPattern)
		}
	} else {
		log.Infof("No matching rule found for %s, forwarding request", r.URL.String())
		logEntry += "    Result: Forwarded (No matching rule)\n"
	}

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusServiceUnavailable)
		log.Errorf("Failed to forward request: %v", err)
		logEntry += fmt.Sprintf("    Error: %v\n", err)
		writeLogEntry(logEntry)
		return
	}
	defer resp.Body.Close()

	log.Infof("Forwarded request to %s with status %d", r.URL.String(), resp.StatusCode)
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Errorf("Failed to copy response body: %v", err)
	}

	duration := time.Since(startTime)
	logEntry += fmt.Sprintf("    Status: %d\n    Duration: %v\n", resp.StatusCode, duration)
	writeLogEntry(logEntry)
}

func initLogger() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func writeLogEntry(entry string) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("Error opening log file: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(entry); err != nil {
		log.Errorf("Error writing to log file: %v", err)
	}
}
