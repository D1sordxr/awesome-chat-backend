package httpFiber

import (
	"awesome-chat/internal/infrastructure/config/http"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"os"
)

type Handler interface {
	RegisterRoutes(router fiber.Router)
}

type Server struct {
	port     string
	handlers []Handler
	Server   *fiber.App
}

func NewServer(cfg *http.Config, handlers ...Handler) *Server {
	app := fiber.New()
	// fiber.Config{DisableStartupMessage: true}

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000",
	}))

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${locals:requestid} ${status} - ${method} ${path} (${latency})\n",
		TimeFormat: "2006-01-02 15:04:05",
		Output:     os.Stdout,
	}))

	// if cfg.Debug {
	app.Use(pprof.New())
	// }

	return &Server{
		port:     cfg.Port,
		handlers: handlers,
		Server:   app,
	}
}
func (s *Server) Start(_ context.Context) error {
	s.setupRoutes()

	if err := s.Server.Listen(":" + s.port); err != nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown(_ context.Context) error {
	if err := s.Server.Shutdown(); err != nil {
		return err
	}
	return nil
}

func (s *Server) setupRoutes() {
	for _, handler := range s.handlers {
		handler.RegisterRoutes(s.Server)
	}
}
