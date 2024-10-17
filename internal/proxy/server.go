package proxy

import (
	"net/http"

	"nexum/internal/config"
	"nexum/internal/logger"
)

type Server struct {
	config *config.Config
	logger *logger.Logger
}

func NewServer(cfg *config.Config, log *logger.Logger) *Server {
	return &Server{
		config: cfg,
		logger: log,
	}
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, newHandler(s.config, s.logger))
}
