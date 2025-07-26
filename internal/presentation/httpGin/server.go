package httpGin

import (
	"awesome-chat/internal/domain/app/ports"
	cfg "awesome-chat/internal/infrastructure/config/http"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler interface {
	RegisterRoutes(router gin.IRouter)
}

type Server struct {
	log      ports.Logger
	handlers []Handler
	engine   *gin.Engine
	server   *http.Server
}

func NewServer(
	log ports.Logger,
	config *cfg.Config,
	handlers ...Handler,
) *Server {
	log.Info("Initializing WS server", "port", config.Port)

	engine := gin.Default()

	return &Server{
		log: log,
		server: &http.Server{
			Addr:              ":" + config.Port,
			Handler:           engine.Handler(),
			ReadHeaderTimeout: config.Timeout,
			ReadTimeout:       config.Timeout,
			WriteTimeout:      config.Timeout,
		},
		engine:   engine,
		handlers: handlers,
	}
}

func (s *Server) Start(_ context.Context) error {
	s.log.Info("Registering WS handlers...")
	for _, handler := range s.handlers {
		handler.RegisterRoutes(s.engine)
	}

	s.log.Info("Starting WS server...", "address", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info("WS server closed gracefully")
			return nil
		}
		s.log.Error("WS server stopped with error", "error", err.Error())
		return err
	}

	s.log.Info("WS server exited unexpectedly")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down WS server...")
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("Failed to gracefully shutdown WS server", "error", err.Error())
		return err
	}
	s.log.Info("WS server shutdown complete")
	return nil
}
