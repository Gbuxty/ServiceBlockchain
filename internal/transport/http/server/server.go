package server

import (
	"BlockchainCurrency/config"
	"BlockchainCurrency/pkg/logger"
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
	cfg        *config.HTTPServer
}

func NewServer(
	cfg *config.HTTPServer,
	logger logger.Logger,
	handler http.Handler,

) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Server) ListenAndServe() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Server) Run() {
	s.logger.Infof("Server listening on port: %d", s.cfg.Port)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Errorf("Listen:", err)
	}

}
