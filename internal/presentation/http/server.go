package http

import (
	"awesome-chat/internal/domain/app/ports"
	cfg "awesome-chat/internal/infrastructure/config/http"
	"context"
	"errors"
	"net/http"
)

type Handler interface {
	RegisterRoutes(mux *http.ServeMux)
}

type Server struct {
	log      ports.Logger
	server   *http.Server
	handlers []Handler
}

func NewServer(
	log ports.Logger,
	config *cfg.Config,
	routes ...Handler,
) *Server {
	log.Info("Initializing HTTP server", "port", config.Port)

	return &Server{
		log: log,
		server: &http.Server{
			Addr:              ":" + config.Port,
			ReadHeaderTimeout: config.Timeout,
			ReadTimeout:       config.Timeout,
			WriteTimeout:      config.Timeout,
		},
		handlers: routes,
	}
}

func (s *Server) Start(_ context.Context) error {
	s.log.Info("Registering HTTP handlers...")
	mux := http.NewServeMux()
	for _, handler := range s.handlers {
		handler.RegisterRoutes(mux)
	}
	s.server.Handler = mux

	s.log.Info("Starting HTTP server...", "address", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info("HTTP server closed gracefully")
			return nil
		}
		s.log.Error("HTTP server stopped with error", "error", err.Error())
		return err
	}

	s.log.Info("HTTP server exited unexpectedly")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down HTTP server...")
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.log.Error("Failed to gracefully shutdown HTTP server", "error", err.Error())
		return err
	}
	s.log.Info("HTTP server shutdown complete")
	return nil
}
