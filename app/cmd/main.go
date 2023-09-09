package main

import (
	"boost-my-skills-bot/config"
	"context"
	"log"
	"sync"

	models "boost-my-skills-bot/internal/models/bot"
	"boost-my-skills-bot/pkg/storage/postgres"
	"os"
	"os/signal"
	"syscall"
	"time"

	"boost-my-skills-bot/internal/bot/repository"
	"boost-my-skills-bot/internal/bot/usecase"

	"boost-my-skills-bot/internal/bot/tgBot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
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
	psqlDB, err := postgres.InitPsqlDB(ctx, cfg)
	if err != nil {
		log.Printf("PostgreSQL error connection: %s", err.Error())
		return
	} else {
		log.Println("PostgreSQL successful connection")
	}
	defer func(psqlDB *sqlx.DB) {
		if err := psqlDB.Close(); err != nil {
			log.Printf("PostgreSQL error close connection: %s", err.Error())
			return
		} else {
			log.Println("PostgreSQL successful close connection")
		}

	}(psqlDB)

	tgbot, err := mapHandler(ctx, cfg, psqlDB)
	if err != nil {
		log.Printf("Error map handler: %s", err.Error())
		return
	}

	if err := tgbot.Run(); err != nil {
		log.Printf("Error bot run: %s", err.Error())
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

}

func mapHandler(ctx context.Context, cfg *config.Config, db *sqlx.DB) (tgBot *tgbot.TgBot, err error) {

	botAPI, err := tgbotapi.NewBotAPI(cfg.TgBot.ApiKey)
	if err != nil {
		return
	}

	stateDirections := sync.Map{}
	userStates := make(map[int64]models.AddQuestionParams)

	// repository
	botRepo := repository.NewBotPGRepo(db)

	// usecase
	botUC := usecase.NewBotUC(cfg, botRepo, botAPI, &stateDirections, userStates)

	// bot
	tgBot = tgbot.NewTgBot(cfg, botUC, botAPI, userStates)

	// map worker
	go func() {
		ticker := time.NewTicker(time.Duration(time.Second * 2))
		for {
			select {
			case <-ticker.C:
				if err = botUC.SyncDirectionsInfo(ctx); err != nil {
					log.Println(err)
				}

			}
		}
	}()

	return tgBot, nil
}
