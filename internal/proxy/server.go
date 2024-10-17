package proxy

import (
	"net/http"
	"nexum/internal/config"
	"nexum/internal/logger"
)

type Server struct {
	config  *config.Config
	logger  *logger.Logger
	handler *Handler
}

func NewServer(cfg *config.Config, log *logger.Logger) *Server {
	server := &Server{
		config: cfg,
		logger: log,
	}
	server.handler = NewHandler(cfg, log)
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
