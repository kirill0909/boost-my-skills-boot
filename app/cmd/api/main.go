package main

import (
	"aifory-pay-admin-bot/config"
	"aifory-pay-admin-bot/internal/server"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"aifory-pay-admin-bot/pkg/storage/postgres"

	"context"

	"go.uber.org/fx"

	"github.com/jmoiron/sqlx"
)

func providePgDB(lifecycle fx.Lifecycle, cfg *config.Config) (*sqlx.DB, error) {
	psqlDB, err := postgres.InitPsqlDB(context.Background(), cfg)
	if err != nil {
		log.Printf("PostgreSQL init error: %s", err.Error())
		return nil, err
	}
	log.Printf("PostgreSQL provided, status: %#v", psqlDB.Stats())

	lifecycle.Append(
		fx.Hook{
			OnStop: func(context.Context) error {
				err := psqlDB.Close()
				if err != nil {
					return err
				}
				log.Println("PostgreSQL closed properly")
				return nil
			},
		},
	)
	return psqlDB, nil
}

func startTraces(cfg *config.Config) (*jaeger.Exporter, *tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.OpenTelemetry.Host)))
	if err != nil {
		return nil, nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.OpenTelemetry.ServiceName),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return exp, tp, nil
}

func runGrpcServer(lifecycle fx.Lifecycle, g *server.GRPCServer) error {
	if err := g.Run(); err != nil {
		return err
	}

	lifecycle.Append(
		fx.Hook{
			OnStop: func(context.Context) error {
				g.GracefulShutdown()
				return nil
			},
		},
	)
	return nil
}

func main() {
	log.Println("Starting server")

	app := fx.New(
		fx.Provide(config.ProvideConfig),
		fx.Provide(providePgDB),
		fx.Provide(server.NewGRPCServer),
		fx.NopLogger,
		fx.Invoke(runGrpcServer),
	)

	jeaeger, traceProvider, err := startTraces(config.C)
	if err != nil {
		log.Printf("Cannot create Jaeger exporter - %s", err.Error())
	}
	log.Println("Jaeger exporter started")
	defer func() {
		err = traceProvider.Shutdown(context.Background())
		if err != nil {
			log.Println(err)
		} else {
			log.Println("OpenTelemetry closed properly")
		}

		err = jeaeger.Shutdown(context.Background())
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Jaeger closed properly")
		}
	}()

	app.Run()

	stopCtx, cancel := context.WithTimeout(context.Background(), config.C.Server.ShutdownTimeout)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Println(err)
	}
}
