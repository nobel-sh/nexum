package proxy

import (
	"net/http"
	// "nexum/internal/config"
	"nexum/internal/logger"
	"nexum/internal/rules"
)

type Server struct {
	// config  *config.Config
	logger  *logger.Logger
	handler *Handler
	rules   *rules.RuleList
}

func NewServer(rules *rules.RuleList, logger *logger.Logger) *Server {
	server := &Server{
		rules:  rules,
		logger: logger,
	}
	server.handler = NewHandler(rules, logger)
	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// if REQUEST method is found create a TCP tunnel between them
	if r.Method == http.MethodConnect {
		s.handler.HandleConnectTunnel(w, r)
	} else {
		s.handler.HandleHTTPRequests(w, r)
	}
}
