package gateway

import (
	"context"
	"errors"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/infrastructure/config/components/cookie"
	httpCfg "awesome-chat/internal/infrastructure/config/components/http"
)

type Registrar func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

type Options struct {
	AllowedOrigins []string
	Cookie         cookie.Config
}

type Server struct {
	log        ports.Logger
	conn       *grpc.ClientConn
	mux        *runtime.ServeMux
	server     *http.Server
	registrars []Registrar
}

func NewServer(
	ctx context.Context,
	log ports.Logger,
	cfg *httpCfg.Config,
	opts Options,
	conn *grpc.ClientConn,
	registrars ...Registrar,
) (*Server, error) {
	mux := newMux(opts.Cookie)

	if err := mux.HandlePath(http.MethodGet, "/healthz", health); err != nil {
		return nil, err
	}

	for _, register := range registrars {
		if err := register(ctx, mux, conn); err != nil {
			log.Error("Failed to register gateway handler", "error", err.Error())

			return nil, err
		}
	}

	return &Server{
		log:  log,
		conn: conn,
		mux:  mux,
		server: &http.Server{
			Addr:              ":" + cfg.Port,
			Handler:           withLogging(withCORS(mux, opts.AllowedOrigins), log),
			ReadHeaderTimeout: cfg.Timeout,
			ReadTimeout:       cfg.Timeout,
			WriteTimeout:      cfg.Timeout,
			IdleTimeout:       cfg.IdleTimeout,
		},
		registrars: registrars,
	}, nil
}

func (s *Server) Start(_ context.Context) error {
	s.log.Info("Starting HTTP gateway", "address", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info("HTTP gateway closed gracefully")

			return nil
		}

		s.log.Error("HTTP gateway stopped with error", "error", err.Error())

		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down HTTP gateway")

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("Failed to shutdown HTTP gateway", "error", err.Error())

		return err
	}

	s.log.Info("HTTP gateway shutdown complete")

	return nil
}

func (s *Server) Handler() http.Handler {
	return s.server.Handler
}

func health(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusOK)
}
