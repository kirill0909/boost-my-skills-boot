package main

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot/repository"
	"boost-my-skills-bot/internal/bot/tgBot"
	"boost-my-skills-bot/internal/bot/usecase"
	models "boost-my-skills-bot/internal/models/bot"
	"boost-my-skills-bot/pkg/storage/postgres"
	"boost-my-skills-bot/pkg/storage/redis"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kirill0909/logger"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	tgbot, err := maping(ctx, cfg, dependencies)
	if err != nil {
		dependencies.Logger.Errorf("Error map handler: %s", err.Error())
		return
	}

	go func() {
		if err := tgbot.Run(); err != nil {
			dependencies.Logger.Errorf("Error bot run: %s", err.Error())
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

}

func maping(ctx context.Context, cfg *config.Config, dep models.Dependencies) (tgBot *tgbot.TgBot, err error) {
	botAPI, err := tgbotapi.NewBotAPI(cfg.TgBot.ApiKey)
	if err != nil {
		return
	}

	// repository
	botPgRepo := repository.NewBotPGRepo(dep.PgDB)
	botRedisRepo := repository.NewBotRedisRepo(dep.Redis, cfg)

	// usecase
	botUC := usecase.NewBotUC(cfg, botPgRepo, botRedisRepo, botAPI, dep.Logger)

	// bot
	tgBot = tgbot.NewTgBot(cfg, botUC, botAPI, dep.Logger)

	go func() {
		botUC.SyncMainKeyboardWorker()
	}()

	return tgBot, nil
}

func initDependencies(ctx context.Context, cfg *config.Config) (models.Dependencies, error) {
	logger := logger.InitLogger()

	pgDB, err := postgres.InitPgDB(ctx, cfg)
	if err != nil {
		return models.Dependencies{}, err
	} else {
		logger.Infof("PostgreSQL successful connection")
	}

	redisDB, _, err := redis.InitRedisClient(cfg)
	if err != nil {
		return models.Dependencies{}, err
	} else {
		logger.Infof("Redis successful connection")
	}

	return models.Dependencies{
		PgDB:   pgDB,
		Redis:  redisDB,
		Logger: logger}, nil
}

func closeDependencies(dep models.Dependencies) error {
	if err := dep.PgDB.Close(); err != nil {
		return errors.Wrap(err, "PostgreSQL error close connection")
	} else {
		dep.Logger.Infof("PostgreSQL successful close connection")
	}

	if err := dep.Redis.Close(); err != nil {
		return errors.Wrap(err, "Redis err close connection")
	} else {
		dep.Logger.Infof("Redis successful close connection")
	}

	return nil
}
