package proxy

import (
	"io"
	"net"
	"net/http"
	"nexum/internal/config"
	"nexum/internal/logger"
	"nexum/internal/rules"
	"nexum/pkg/httputil"
	"time"
)

type Handler struct {
	config *config.Config
	logger *logger.Logger
}

func NewHandler(cfg *config.Config, log *logger.Logger) *Handler {
	return &Handler{
		config: cfg,
		logger: log,
	}
}

func (h *Handler) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	h.logger.Info("Received HTTP request %s %s from %s", r.Method, r.URL.String(), r.RemoteAddr)

	rule := rules.MatchRule(h.config.Rules, r.URL.String())
	if rule != nil {
		switch rule.Action {
		case "block":
			http.Error(w, "Access denied", http.StatusForbidden)
			h.logger.Info("Blocked request %s based on rule %s", r.URL.String(), rule.URLPattern)
			return
		case "modify":
			rules.ApplyModifications(r, rule.Modifications)
			h.logger.Info("Modified request %s based on rule %s", r.URL.String(), rule.URLPattern)
		case "allow":
			h.logger.Info("Allowed request %s based on rule %s", r.URL.String(), rule.URLPattern)
		}
	} else {
		h.logger.Info("No matching rule found for %s, forwarding request", r.URL.String())
	}

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusServiceUnavailable)
		h.logger.Error("Failed to forward request: %v", err)
		return
	}
	defer resp.Body.Close()

	h.logger.Info("Forwarded request to %s with status %d", r.URL.String(), resp.StatusCode)
	httputil.CopyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	duration := time.Since(startTime)
	h.logger.Info("HTTP request processed in %v", duration)
}

func (h *Handler) HandleHTTPS(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	h.logger.Info("Received HTTPS CONNECT request for %s from %s", r.Host, r.RemoteAddr)

	targetConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		h.logger.Error("Failed to connect to %s: %v", r.Host, err)
		return
	}

	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		h.logger.Error("Hijacking not supported")
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		h.logger.Error("Hijacking failed: %v", err)
		return
	}

	h.logger.Info("HTTPS connection established between %s and %s", r.RemoteAddr, r.Host)

	go h.transfer(targetConn, clientConn)
	go h.transfer(clientConn, targetConn)

	duration := time.Since(startTime)
	h.logger.Info("HTTPS CONNECT request processed in %v", duration)
}

func (h *Handler) transfer(destination io.WriteCloser, source io.ReadCloser) {
	_, err := io.Copy(destination, source)
	destination.Close()
	source.Close()

	if err != nil {
		h.logger.Error("Error during data transfer: %v", err)
	}
}
