package proxy

import (
	"io"
	"net/http"
	"time"

	"nexum/internal/config"
	"nexum/internal/logger"
	"nexum/internal/rules"
	"nexum/pkg/httputil"
)

type handler struct {
	config *config.Config
	logger *logger.Logger
}

func newHandler(cfg *config.Config, log *logger.Logger) *handler {
	return &handler{
		config: cfg,
		logger: log,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	h.logger.Info("Received request %s %s from %s", r.Method, r.URL.String(), r.RemoteAddr)

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
	if _, err := io.Copy(w, resp.Body); err != nil {
		h.logger.Error("Failed to copy response body: %v", err)
	}

	duration := time.Since(startTime)
	h.logger.Info("Request processed in %v", duration)
}
