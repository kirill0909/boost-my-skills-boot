package server

import (
	"aifory-pay-admin-bot/internal/bot/repository"
	"aifory-pay-admin-bot/internal/bot/usecase"

	bot2 "aifory-pay-admin-bot/internal/bot/tgBot"

	handlerBot "aifory-pay-admin-bot/internal/bot/handlers"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	pb "gitlab.axarea.ru/main/aiforypay/package/admin-bot-proto"
)

func (g *GRPCServer) MapHandlers() error {
	botAPI, err := tgbotapi.NewBotAPI(g.cfg.Server.TGToken)
	if err != nil {
		log.Printf("Bot create error: %s", err.Error())
		return err
	}

	botPGRepo := repository.NewBotRepository(g.pgDB)

	botUc := usecase.NewBotUsecase(botPGRepo, g.cfg, botAPI)
	bot := bot2.NewTgBot(g.cfg, botPGRepo, botUc, botAPI)

	go func() {
		log.Print("Bot is running")
		if err := bot.Run(); err != nil {
			log.Fatalf("Error bot.Run: %s", err.Error())
		}
	}()

	botHandler := handlerBot.NewBotHandler(botUc)
	pb.RegisterAdminBotServer(g.srv, botHandler)
	return nil
}
