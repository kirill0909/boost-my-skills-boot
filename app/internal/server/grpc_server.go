package server

import (
	"aifory-pay-admin-bot/config"

	"log"
	"net"

	"github.com/gofiber/fiber/v2"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jmoiron/sqlx"
	grpc_zerolog "github.com/philip-bui/grpc-zerolog"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"os"

	"github.com/rs/zerolog"
)

var (
	zerologger = zerolog.New(os.Stdout).With().Timestamp().Logger()
)

type GRPCServer struct {
	srv  *grpc.Server
	pgDB *sqlx.DB
	cfg  *config.Config
}

func NewGRPCServer(cfg *config.Config, pgDB *sqlx.DB) *GRPCServer {
	return &GRPCServer{
		cfg:  cfg,
		pgDB: pgDB,
	}
}

func (g *GRPCServer) Run() error {
	l, err := net.Listen("tcp", g.cfg.Server.Host+":"+g.cfg.Server.GRPCPort)
	if err != nil {
		return err
	}

	g.srv = grpc.NewServer(
		// grpc.KeepaliveParams(
		// 	keepalive.ServerParameters{
		// 		MaxConnectionIdle: s.cfg.PaymentServer.MaxConnectionIdle * time.Minute,
		// 		Timeout:           s.cfg.PaymentServer.Timeout * time.Second,
		// 		MaxConnectionAge:  s.cfg.PaymentServer.MaxConnectionAge * time.Minute,
		// 		Time:              s.cfg.PaymentServer.Timeout * time.Minute,
		// 	},
		// ),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				otelgrpc.UnaryServerInterceptor(),
				grpc_zerolog.NewUnaryServerInterceptorWithLogger(&zerologger),
			),
		),
	)

	if err := g.MapHandlers(); err != nil {
		return errors.Wrapf(err, "cannot map handlers")
	}

	go func() {
		f := fiber.New(fiber.Config{DisableStartupMessage: true})
		f.Get("/health_check", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		if err := f.Listen(g.cfg.Server.Host + ":" + g.cfg.Server.HTTPPort); err != nil {
			log.Fatalf("HTTP Server error: %s", err.Error())
		}
	}()

	go func() {
		log.Printf("GRPC Server started on %s:%s", g.cfg.Server.Host, g.cfg.Server.GRPCPort)
		if err := g.srv.Serve(l); err != nil {
			log.Fatalf("GRPC Server error: %s", err.Error())
		}
	}()
	return nil
}

func (g *GRPCServer) GracefulShutdown() {
	g.srv.GracefulStop()
	log.Println("GRPC Server closed properly")
}
