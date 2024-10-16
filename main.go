package main

import (
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func InitLogger() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	InitLogger()
	listenAddr := ":8080"

	log.Infof("Starting proxy server on %s...", listenAddr)
	if err := http.ListenAndServe(listenAddr, http.HandlerFunc(handleRequest)); err != nil {
		log.Errorf("Failed to start server: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Infof("Received request %s %s", r.Method, r.URL.String())

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusServiceUnavailable)
		log.Errorf("Failed to forward request: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Infof("Forwarded request to %s with status %d", r.URL.String(), resp.StatusCode)

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Errorf("Failed to copy response body: %v", err)
	}
}

func copyHeader(dst, src http.Header) {
	for k, v := range src {
		dst[k] = v
	}
}
