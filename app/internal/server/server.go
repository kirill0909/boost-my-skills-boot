package server

import (
	"fmt"
	"net"

	pb "boost-my-skills-bot/app/pkg/proto/github.com/kirill0909/boost-my-skills-boot/app/pkg/proto/boost_bot_proto"

	statisticsAdapter "boost-my-skills-bot/app/internal/statistics/adapter"

	"github.com/gofiber/fiber/v2"
	"github.com/kirill0909/logger"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Server struct {
	HTTP HTTP
	GRPC GRPC
	log  *logger.Logger
}

type GRPC struct {
	srv         *grpc.Server
	statAdapter *statisticsAdapter.Statistics
	host        string
	port        string
}

type HTTP struct {
	app  *fiber.App
	host string
	port string
}

func NewServer(HTTPHost, HTTPport, GRPCHost, GRPCPort string, logger *logger.Logger, statAdapter *statisticsAdapter.Statistics) *Server {
	return &Server{
		HTTP: HTTP{app: fiber.New(), host: HTTPHost, port: HTTPport},
		GRPC: GRPC{srv: grpc.NewServer(), statAdapter: statAdapter, host: GRPCHost, port: GRPCPort},
		log:  logger,
	}
}

func (s *Server) RunGRPC() error {
	pb.RegisterStatisticsServer(s.GRPC.srv, s.GRPC.statAdapter)

	addr := fmt.Sprintf("%s:%s", s.GRPC.host, s.GRPC.port)
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

	addr := fmt.Sprintf("%s:%s", s.HTTP.host, s.HTTP.port)
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
