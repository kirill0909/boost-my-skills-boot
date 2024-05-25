package server

import (
	"fmt"
	"net"

	pb "boost-my-skills-bot/app/pkg/proto/github.com/kirill0909/boost-my-skills-boot/app/pkg/proto/boost_bot_proto"

	"github.com/gofiber/fiber/v2"
	"github.com/kirill0909/logger"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Server struct {
	HTTP HTTP
	GRPC GRPC
	host string
	log  *logger.Logger
}

type GRPC struct {
	srv  *grpc.Server
	port string
}

type HTTP struct {
	app  *fiber.App
	port string
}

func NewServer(host, HTTPport, GRPCPort string, logger *logger.Logger) *Server {
	return &Server{
		HTTP: HTTP{app: fiber.New(), port: HTTPport},
		GRPC: GRPC{srv: grpc.NewServer(), port: GRPCPort},
		host: host,
		log:  logger,
	}
}

func (s *Server) RunGRPC() error {
	pb.RegisterBostBotServer(s.GRPC.srv, pb.UnimplementedBostBotServer{})

	addr := fmt.Sprintf("%s:%s", s.host, s.GRPC.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "Server.RunHTTP().Listen()")
	}

	if err := s.GRPC.srv.Serve(listener); err != nil {
		return errors.Wrap(err, "Server.RunGRPC.Server()")
	}

	return nil
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

func (s *Server) ShutdownGRPC() {
	s.GRPC.srv.GracefulStop()
}
