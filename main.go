package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	listenAddr := ":8080"

	log.Printf("INFO: Starting proxy server on %s...", listenAddr)
	if err := http.ListenAndServe(listenAddr, http.HandlerFunc(handleRequest)); err != nil {
		log.Fatalf("ERROR: Failed to start server: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: Received request %s %s", r.Method, r.URL.String())

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusServiceUnavailable)
		log.Printf("ERROR: Failed to forward request: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("INFO: Forwarded request to %s with status %d", r.URL.String(), resp.StatusCode)

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("ERROR: Failed to copy response body: %v", err)
	}
}

func copyHeader(dst, src http.Header) {
	for k, v := range src {
		dst[k] = v
	}
}
