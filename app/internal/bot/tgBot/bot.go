package tgbot

import (
	"boost-my-skills-bot/app/config"
	"boost-my-skills-bot/app/internal/bot"
	"boost-my-skills-bot/app/internal/bot/models"
	"boost-my-skills-bot/app/pkg/utils"
	"context"
	"encoding/json"

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
					t.sendErrorMessage(update.Message.Chat.ID, "Account activation error. Contact @kirillkorunov to get uuid for account actionvation")
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
			case utils.PrintQuestionsCommand:
				params := models.HandlePrintQuestionsCommandParams{ChatID: update.Message.Chat.ID}
				if err := t.tgUC.HandlePrintQuestionsCommand(ctx, params); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "print info error")
					continue
				}
				continue
			}

			// handle entered text
			switch {
			case statusID == utils.AwaitingDirectionNameStatus: // execute when user enter direction name
				params := models.CreateDirectionParams{ChatID: update.Message.Chat.ID, DirectionName: update.Message.Text}
				if err := t.tgUC.CreateDirection(ctx, params); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "internal server error")
					continue
				}
				continue
			case statusID == utils.AwaitingQuestionStatus: // execute when user enter question
				params := models.HandleAwaitingQuestionParams{ChatID: update.Message.Chat.ID, Question: update.Message.Text}
				if err := t.tgUC.HandleAwaitingQuestion(ctx, params); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "internal server error")
					continue
				}
				continue
			case statusID == utils.AwaitingAnswerStatus: // execute when user enter answer
				params := models.HandleAwaitingAnswerParams{ChatID: update.Message.Chat.ID, Answer: update.Message.Text}
				if err := t.tgUC.HandleAwaitingAnswer(ctx, params); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "internal server error")
					continue
				}
				continue
			default:
				t.sendMessage(update.Message.From.ID, "use keyboard to interact with bot")
				continue
			}
		}

		// handle callbacks
		if update.CallbackQuery != nil {
			callbackInfo, err := t.extractCallbackInfo(update.CallbackData())
			if err != nil {
				t.log.Errorf(err.Error())
				t.sendErrorMessage(update.Message.Chat.ID, "internal server error")
				continue
			}
			switch {
			case callbackInfo.CallbackType == utils.AwaitingParentDirectionCallbackType: // executes when user tap on direction name button
				parentDirectionParams := models.SetParentDirectionParams{ChatID: update.CallbackQuery.From.ID, ParentDirectionID: callbackInfo.DirectionID}
				if err := t.tgUC.SetParentDirection(ctx, parentDirectionParams); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "internal server error")
					continue
				}

				createDirectionCommandParams := models.HandleCreateDirectionCommandParams{
					ChatID:            update.CallbackQuery.From.ID,
					ParentDirectionID: callbackInfo.DirectionID,
				}
				if err := t.tgUC.HandleCreateDirectionCommand(ctx, createDirectionCommandParams); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "internal server error")
					continue
				}
				continue
			case callbackInfo.CallbackType == utils.AwaitingAddInfoDirectionCallbackType:
				params := models.HandleAddInfoCommandParams{
					ChatID:            update.CallbackQuery.From.ID,
					ParentDirectionID: callbackInfo.DirectionID,
				}
				if err := t.tgUC.HandleAddInfoCommand(ctx, params); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "internal server error")
					continue
				}
				continue
			case callbackInfo.CallbackType == utils.AwaitingPrintQuestionsCallbackType:
				params := models.HandlePrintQuestionsCommandParams{
					ChatID:            update.CallbackQuery.From.ID,
					ParentDirectionID: callbackInfo.DirectionID,
				}
				if err := t.tgUC.HandlePrintQuestionsCommand(ctx, params); err != nil {
					t.log.Errorf(err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, "internal server error")
					continue
				}
				continue
			case callbackInfo.CallbackType == utils.AwaitingInfoActionsCallbackType:
				params := models.HandleAwaitingPrintAnswerParams{
					ChatID:    update.CallbackQuery.From.ID,
					InfoID:    callbackInfo.InfoID,
					MessageID: update.CallbackQuery.Message.MessageID,
				}
				if err := t.tgUC.HandleAwaitingPrintAnswer(ctx, params); err != nil {
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

func (t *TgBot) extractCallbackInfo(callbackData string) (models.CallbackInfo, error) {
	var callbackInfo models.CallbackInfo
	if err := json.Unmarshal([]byte(callbackData), &callbackInfo); err != nil {
		return models.CallbackInfo{}, err
	}

	return callbackInfo, nil
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
