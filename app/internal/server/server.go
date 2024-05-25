package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kirill0909/logger"
	"github.com/pkg/errors"
)

type Server struct {
	HTTP HTTP
	host string
	log  *logger.Logger
}

type HTTP struct {
	app  *fiber.App
	port string
}

func NewServer(host, HTTPport string, logger *logger.Logger) *Server {
	return &Server{
		HTTP: HTTP{app: fiber.New(), port: HTTPport},
		host: host,
		log:  logger,
	}
}

func (s *Server) RunHTTP() error {

	s.HTTP.app.Get("/ping", func(c *fiber.Ctx) error {
		s.log.Infof("path: %s | status: %d", string(c.Context().RequestURI()), fiber.StatusOK)
		return c.SendString("pong")
	})

	addr := fmt.Sprintf("%s:%s", s.host, s.HTTP.port)
	if err := s.HTTP.app.Listen(addr); err != nil {
		return errors.Wrap(err, "Server.RunHTTP().Listen()")
	}

	return nil
}

func (s *Server) ShutdownHTTP() error {
	if err := s.HTTP.app.Shutdown(); err != nil {
		return errors.Wrap(err, "Server.ShutdownHTTP().Shutdown()")
	}

	return nil
}
