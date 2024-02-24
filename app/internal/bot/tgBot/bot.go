package tgbot

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	"boost-my-skills-bot/internal/bot/models"
	"boost-my-skills-bot/pkg/utils"
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
		ctx := context.Background()
		statusID, err := t.tgUC.GetAwaitingStatus(ctx, update.FromChat().ID)
		if err != nil {
			t.log.Errorf(err.Error())
			t.sendErrorMessage(update.Message.Chat.ID, "internal server error")
			continue
		}

		// handle message
		if update.Message != nil {
			switch update.Message.Command() {
			case utils.StartCommand:
				params := models.HandleStartCommandParams{
					Text: update.Message.Text, ChatID: update.Message.Chat.ID, TgName: update.Message.Chat.UserName}
				if err := t.tgUC.HandleStartCommand(ctx, params); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "account activation error")
					continue
				}
				continue
			case utils.CreateDirectionCommand:
				params := models.HandleCreateDirectionCommandParams{
					Text: update.Message.Text, ChatID: update.Message.Chat.ID}
				if err := t.tgUC.HandleCreateDirectionCommand(ctx, params); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "create direction error")
					continue
				}
				continue
			case utils.AddInfoCommand:
				params := models.HandleAddInfoCommandParams{ChatID: update.Message.Chat.ID}
				if err := t.tgUC.HandleAddInfoCommand(ctx, params); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "add info error")
					continue
				}
				continue
			}

			// handle entered text
			switch {
			case statusID == utils.AwaitingDirectionNameStatus || statusID == utils.AwaitingParentDirecitonStatus: // execute when user enter direction name
				params := models.CreateDirectionParams{ChatID: update.Message.Chat.ID, DirectionName: update.Message.Text}
				if err := t.tgUC.CreateDirection(ctx, params); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "internal server error")
					continue
				}
				continue
			case statusID == utils.AwaitingQuestion: // execute when user enter question
				t.log.Infof("Yout question is: %s", update.Message.Text)
				continue
			default:
				t.sendMessage(update.Message.From.ID, "use keyboard to interact with bot")
				continue
			}
		}

		// handle callbacks
		if update.CallbackQuery != nil {
			switch {
			case statusID == utils.AwaitingParentDirecitonStatus: // executes when user tap on direction name button
				parentDirectionParams := models.SetParentDirectionParams{ChatID: update.CallbackQuery.From.ID, CallbackData: update.CallbackData()}
				if err := t.tgUC.SetParentDirection(ctx, parentDirectionParams); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "internal server error")
					continue
				}

				createDirectionCommandParams := models.HandleCreateDirectionCommandParams{
					ChatID:       update.CallbackQuery.From.ID,
					CallbackData: update.CallbackData()}
				if err := t.tgUC.HandleCreateDirectionCommand(ctx, createDirectionCommandParams); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "internal server error")
					continue
				}
				continue
			case statusID == utils.AwaitingAddInfoDirection:
				params := models.HandleAddInfoCommandParams{
					ChatID:       update.CallbackQuery.From.ID,
					CallbackData: update.CallbackData()}
				if err := t.tgUC.HandleAddInfoCommand(ctx, params); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "internal server error")
					continue
				}
				continue
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

func (t *TgBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := t.BotAPI.Send(msg); err != nil {
		t.log.Errorf(err.Error())
	}
}
