package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kirill0909/logger"
	"github.com/pkg/errors"
)

type Server struct {
	srv  *fiber.App
	host string
	port string
	log  *logger.Logger
}

func NewServer(host, port string, logger *logger.Logger) *Server {
	return &Server{srv: fiber.New(), host: host, port: port, log: logger}
}

func (s *Server) Run() error {

	s.srv.Get("/ping", func(c *fiber.Ctx) error {
		s.log.Infof("path: %s | status: %d", string(c.Context().RequestURI()), fiber.StatusOK)
		return c.SendString("pong")
	})

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	if err := s.srv.Listen(addr); err != nil {
		return errors.Wrap(err, "Server.Run().Listen()")
	}

	return nil
}

func (s *Server) Shutdown() error {
	if err := s.srv.Shutdown(); err != nil {
		return errors.Wrap(err, "Server.Shutdown().Shutdown()")
	}

	return nil
}
