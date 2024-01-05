package main

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot/repository"
	"boost-my-skills-bot/internal/bot/tgBot"
	"boost-my-skills-bot/internal/bot/usecase"
	models "boost-my-skills-bot/internal/models/bot"
	"boost-my-skills-bot/pkg/storage/postgres"
	"boost-my-skills-bot/pkg/storage/rabbit"
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

	tgbot, err := mapHandler(ctx, cfg, dependencies)
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

func mapHandler(ctx context.Context, cfg *config.Config, dep models.Dependencies) (tgBot *tgbot.TgBot, err error) {
	botAPI, err := tgbotapi.NewBotAPI(cfg.TgBot.ApiKey)
	if err != nil {
		return
	}

	// repository
	botRepo := repository.NewBotPGRepo(dep.PgDB)

	// usecase
	botUC := usecase.NewBotUC(cfg, botRepo, dep.RabbitMQ, botAPI, dep.Logger)

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

	rabbitMQ, err := rabbit.InitRabbit(cfg)
	if err != nil {
		return models.Dependencies{}, err
	} else {
		logger.Infof("RabbitMQ successfil initialization")
	}

	return models.Dependencies{PgDB: pgDB, Logger: logger, RabbitMQ: rabbitMQ}, nil
}

func closeDependencies(dep models.Dependencies) error {
	if err := dep.PgDB.Close(); err != nil {
		return errors.Wrap(err, "PostgreSQL error close connection")
	} else {
		dep.Logger.Infof("PostgreSQL successful close connection")
	}

	if err := dep.RabbitMQ.Chann.Close(); err != nil {
		return errors.Wrap(err, "RabbitMQ error close chann")
	} else {
		dep.Logger.Infof("RabbitMQ successful close chann")
	}

	if err := dep.RabbitMQ.Conn.Close(); err != nil {
		return errors.Wrap(err, "RabbitMQ error close connection")
	} else {
		dep.Logger.Infof("RabbitMQ successful close connection")
	}

	return nil
}
