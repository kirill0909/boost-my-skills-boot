package main

import (
	"boost-my-skills-bot/app/config"
	"boost-my-skills-bot/app/internal/bot/repository"
	tgbot "boost-my-skills-bot/app/internal/bot/tgBot"
	"boost-my-skills-bot/app/internal/bot/usecase"
	"boost-my-skills-bot/app/internal/models"
	"boost-my-skills-bot/app/internal/server"
	"boost-my-skills-bot/app/pkg/storage/postgres"
	"boost-my-skills-bot/app/pkg/storage/redis"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kirill0909/logger"
	"github.com/pkg/errors"
)

func main() {
	viper, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.ParseConfig(viper)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config loaded")

	ctx := context.Background()
	dependencies, err := initDependencies(ctx, cfg)
	if err != nil {
		log.Println(err.Error())
	}

	defer func(dependencies models.Dependencies) {
		if err := closeDependencies(dependencies); err != nil {
			dependencies.Logger.Errorf(err.Error())
		}
	}(dependencies)

	srv, err := maping(cfg, dependencies)
	if err != nil {
		dependencies.Logger.Errorf("error map handler: %s", err.Error())
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	if err := srv.ShutdownHTTP(); err != nil {
		dependencies.Logger.Errorf("error server shudown: %s", err.Error())
		return
	}

	srv.ShutdownGRPC()

	dependencies.Logger.Infof("graceful shutdown server")
}

func maping(cfg *config.Config, dep models.Dependencies) (*server.Server, error) {
	botAPI, err := tgbotapi.NewBotAPI(cfg.TgBot.ApiKey)
	if err != nil {
		return nil, errors.Wrap(err, "maping.NewBotAPI()")
	}

	// repository
	botPgRepo := repository.NewBotPGRepo(dep.PgDB)
	botRedisRepo := repository.NewBotRedisRepo(dep.Redis, cfg)

	// usecase
	botUC := usecase.NewBotUC(cfg, botPgRepo, botRedisRepo, dep.RedisPubSub, botAPI, dep.Logger)

	// bot(adapter)
	bot := tgbot.NewTgBot(cfg, botUC, botAPI, dep.Logger)

	go func(bot *tgbot.TgBot) {
		if err := bot.Run(); err != nil {
			dep.Logger.Fatalf("unable to run bot: %s", err.Error())
		}
	}(bot)

	srv := server.NewServer(cfg.Server.Host, cfg.Server.HTTP.Port, cfg.Server.GRPC.Port, dep.Logger)
	go func(s server.HTTP) {
		if err := srv.RunHTTP(); err != nil {
			dep.Logger.Fatalf("unable to run http server: %s", err.Error())
		}
	}(srv.HTTP)

	go func(s server.GRPC) {
		if err := srv.RunGRPC(); err != nil {
			dep.Logger.Fatalf("unable to run grpc server: %s", err.Error())
		}
	}(srv.GRPC)

	// workers
	go botUC.SyncMainKeyboardWorker()
	go botUC.ListenExpiredMessageWorker()

	return srv, nil
}

func initDependencies(ctx context.Context, cfg *config.Config) (models.Dependencies, error) {
	logger := logger.InitLogger()

	pgDB, err := postgres.InitPgDB(ctx, cfg)
	if err != nil {
		return models.Dependencies{}, err
	} else {
		logger.Infof("PostgreSQL successful connection")
	}

	redisDB, redisPubSub, err := redis.InitRedisClient(cfg)
	if err != nil {
		return models.Dependencies{}, err
	} else {
		logger.Infof("Redis successful connection")
	}

	return models.Dependencies{
		PgDB:        pgDB,
		Redis:       redisDB,
		RedisPubSub: redisPubSub,
		Logger:      logger}, nil
}

func closeDependencies(dep models.Dependencies) error {
	if err := dep.PgDB.Close(); err != nil {
		return errors.Wrap(err, "PostgreSQL error close connection")
	} else {
		dep.Logger.Infof("PostgreSQL successful close connection")
	}

	if err := dep.RedisPubSub.Close(); err != nil {
		return errors.Wrap(err, "Redis error close PubSub connection")
	} else {
		dep.Logger.Infof("Redis successful close PubSub connection")
	}

	if err := dep.Redis.Close(); err != nil {
		return errors.Wrap(err, "Redis err close connection")
	} else {
		dep.Logger.Infof("Redis successful close connection")
	}

	return nil
}
