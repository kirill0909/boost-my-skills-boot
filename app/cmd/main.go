package main

import (
	"boost-my-skills-bot/app/config"
	adapterBot "boost-my-skills-bot/app/internal/bot/adapter"
	repositoryBot "boost-my-skills-bot/app/internal/bot/repository"
	useCaseBot "boost-my-skills-bot/app/internal/bot/usecase"
	"boost-my-skills-bot/app/internal/models"
	"boost-my-skills-bot/app/internal/server"
	adapterStatistics "boost-my-skills-bot/app/internal/statistics/adapter"
	repositoryStatistics "boost-my-skills-bot/app/internal/statistics/repository"
	useCaseStatistics "boost-my-skills-bot/app/internal/statistics/usecase"
	"boost-my-skills-bot/app/pkg/logger"
	"boost-my-skills-bot/app/pkg/storage/postgres"
	"boost-my-skills-bot/app/pkg/storage/redis"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func main() {
	log := logger.NewLogger()

	viper, err := config.LoadConfig()
	if err != nil {
		log.Error("main.LoadConfig()", "error", err.Error())
		return
	}

	cfg, err := config.ParseConfig(viper)
	if err != nil {
		log.Error("main.ParseConfig()", "error", err.Error())
		return
	}
	log.Info("main()", "info", "config loaded")

	ctx := context.Background()
	dependencies, err := initDependencies(ctx, cfg, log)
	if err != nil {
		log.Error("main.initDependencies()", "error", err.Error())
		return
	}

	defer func(dependencies models.Dependencies) {
		if err := closeDependencies(dependencies, log); err != nil {
			log.Error("main.closeDependencies()", "error", err.Error())
			return
		}
	}(dependencies)

	srv, err := maping(cfg, dependencies, log)
	if err != nil {
		log.Error("main.maping()", "error", err.Error())
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	if err := srv.ShutdownHTTP(); err != nil {
		log.Error("main.ShutdownHTTP()", "error", err.Error())
		return
	}

	srv.ShutdownGRPC()

	log.Info("main()", "info", "graceful shutdown server")
}

func maping(cfg *config.Config, dep models.Dependencies, log *slog.Logger) (*server.Server, error) {
	botAPI, err := tgbotapi.NewBotAPI(cfg.TgBot.ApiKey)
	if err != nil {
		return nil, errors.Wrap(err, "maping.NewBotAPI()")
	}

	// repository
	botPgRepo := repositoryBot.NewBotPGRepo(dep.PgDB)
	botRedisRepo := repositoryBot.NewBotRedisRepo(dep.Redis, cfg)
	statisticsPgRepo := repositoryStatistics.NewStatisticsPgRepo(dep.PgDB)

	// usecase
	botUC := useCaseBot.NewBotUC(cfg, botPgRepo, botRedisRepo, dep.RedisPubSub, botAPI, log)
	statisticsUC := useCaseStatistics.NewStatisticsUsecase(statisticsPgRepo)

	// adapter
	botAdapter := adapterBot.NewTgBot(cfg, botUC, botAPI, log)
	statisticsAdapter := adapterStatistics.NewStatistics(statisticsUC, log)

	go func(bot *adapterBot.TgBot) {
		if err := bot.Run(); err != nil {
			log.Error("maping.Run(). unable to run bot", "error", err.Error())
			return
		}
	}(botAdapter)

	srv := server.NewServer(cfg.Server.HTTP.Host, cfg.Server.HTTP.Port, cfg.Server.GRPC.Host, cfg.Server.GRPC.Port, log, statisticsAdapter, cfg.GRPCApiKey)
	go func(s server.HTTP) {
		if err := srv.RunHTTP(); err != nil {
			log.Error("mapgin.RunHTTP(). unable to run http server", "error", err.Error())
			return
		}
	}(srv.HTTP)

	go func(s server.GRPC) {
		if err := srv.RunGRPC(); err != nil {
			log.Error("mapgin.RunGRPC(). unable to run grpc server", "error", err.Error())
			return
		}
	}(srv.GRPC)

	// workers
	go botUC.SyncMainKeyboardWorker()
	go botUC.ListenExpiredMessageWorker()

	return srv, nil
}

func initDependencies(ctx context.Context, cfg *config.Config, log *slog.Logger) (models.Dependencies, error) {
	pgDB, err := postgres.InitPgDB(ctx, cfg)
	if err != nil {
		return models.Dependencies{}, err
	} else {
		log.Info("initDependencies.InitPgDB()", "info", "PostgreSQL successful connection")
	}

	redisDB, redisPubSub, err := redis.InitRedisClient(cfg)
	if err != nil {
		return models.Dependencies{}, err
	} else {
		log.Info("initDependencies.InitRedisClient()", "info", "Redis successful connection")
	}

	return models.Dependencies{
		PgDB:        pgDB,
		Redis:       redisDB,
		RedisPubSub: redisPubSub}, nil
}

func closeDependencies(dep models.Dependencies, log *slog.Logger) error {
	if err := dep.PgDB.Close(); err != nil {
		return errors.Wrap(err, "PostgreSQL error close connection")
	} else {
		log.Info("closeDependencies.Close()", "info", "PostgreSQL successful close connection")
	}

	if err := dep.RedisPubSub.Close(); err != nil {
		return errors.Wrap(err, "Redis error close PubSub connection")
	} else {
		log.Info("closeDependencies.Close()", "info", "Redis successful close PubSub connection")
	}

	if err := dep.Redis.Close(); err != nil {
		return errors.Wrap(err, "Redis err close connection")
	} else {
		log.Info("closeDependencies.Close()", "info", "Redis successful close connection")
	}

	return nil
}
