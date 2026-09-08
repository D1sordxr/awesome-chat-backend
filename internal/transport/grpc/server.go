package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/infrastructure/config/components/grpcserver"
)

type Registrar func(server grpc.ServiceRegistrar)

type Server struct {
	log        ports.Logger
	address    string
	server     *grpc.Server
	registrars []Registrar
}

func NewServer(
	log ports.Logger,
	cfg grpcserver.Config,
	interceptors []grpc.UnaryServerInterceptor,
	registrars ...Registrar,
) *Server {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))

	for _, register := range registrars {
		register(server)
	}

	return &Server{
		log:        log,
		address:    cfg.Address(),
		server:     server,
		registrars: registrars,
	}
}

func (s *Server) Start(_ context.Context) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		s.log.Error("Failed to listen", "address", s.address, "error", err.Error())

		return err
	}

	s.log.Info("Starting gRPC server", "address", s.address)

	if err = s.server.Serve(listener); err != nil {
		s.log.Error("gRPC server stopped with error", "error", err.Error())

		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down gRPC server")

	stopped := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		s.log.Info("gRPC server shutdown complete")
	case <-ctx.Done():
		s.server.Stop()
		s.log.Warn("gRPC server stopped forcefully")
	}

	return nil
}
