package tgbot

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kirill0909/logger"
)

type TgBot struct {
	BotAPI *tgbotapi.BotAPI
	cfg    *config.Config
	tgUC   bot.Usecase
	log    *logger.Logger
}

func NewTgBot(
	cfg *config.Config,
	usecase bot.Usecase,
	botAPI *tgbotapi.BotAPI,
	log *logger.Logger,
) *TgBot {
	return &TgBot{
		cfg:    cfg,
		BotAPI: botAPI,
		tgUC:   usecase,
		log:    log,
	}
}

func (t *TgBot) Run() error {
	t.log.Infof("Authorized on account %s", t.BotAPI.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := t.BotAPI.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			ctx := context.Background()
			switch update.Message.Command() {
			case startCommand:
				params := models.HandleStartCommandParams{
					Text: update.Message.Text, ChatID: update.Message.Chat.ID, TgName: update.Message.Chat.UserName}
				if err := t.handleStartCommand(ctx, params); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "account activation error")
					continue
				}
			case createDirection:
				params := models.HandleCreateDirectionCommandParams{
					Text: update.Message.Text, ChatID: update.Message.Chat.ID, TgName: update.Message.Chat.UserName}
				if err := t.handleCreateDirectionCommand(ctx, params); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "create direction error")
					continue
				}
			}
		}
	}

	return nil
}

func (t *TgBot) sendErrorMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := t.BotAPI.Send(msg)
	if err != nil {
		t.log.Errorf(err.Error())
	}
}
