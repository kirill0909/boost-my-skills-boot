package main

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot/repository"
	"boost-my-skills-bot/internal/bot/tgBot"
	"boost-my-skills-bot/internal/bot/usecase"
	models "boost-my-skills-bot/internal/models/bot"
	"boost-my-skills-bot/pkg/logger"
	"boost-my-skills-bot/pkg/storage/postgres"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	// "github.com/jmoiron/sqlx"
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
			log.Println(err)
		}
	}(dependencies)

	tgbot, err := mapHandler(ctx, cfg, dependencies)
	if err != nil {
		log.Printf("Error map handler: %s", err.Error())
		return
	}

	go func() {
		if err := tgbot.Run(); err != nil {
			log.Printf("Error bot run: %s", err.Error())
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
	botUC := usecase.NewBotUC(cfg, botRepo, botAPI, dep.Logger)

	// bot
	tgBot = tgbot.NewTgBot(cfg, botUC, botAPI, dep.Logger)

	go func() {
		botUC.SyncMainKeyboardWorker()
	}()

	return tgBot, nil
}

func initDependencies(ctx context.Context, cfg *config.Config) (models.Dependencies, error) {
	pgDB, err := postgres.InitPgDB(ctx, cfg)
	if err != nil {
		return models.Dependencies{}, err
	} else {
		log.Println("PostgreSQL successful connection")
	}

	logger := logger.InitLogger()

	return models.Dependencies{PgDB: pgDB, Logger: logger}, nil
}

func closeDependencies(dep models.Dependencies) error {
	if err := dep.PgDB.Close(); err != nil {
		return errors.Wrap(err, "PostgreSQL error close connection")
	} else {
		log.Println("PostgreSQL successful close connection")
	}

	return nil
}
