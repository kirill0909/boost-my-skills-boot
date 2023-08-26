package tgbot

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	"context"
	"fmt"
	"log"
	"strings"

	models "boost-my-skills-bot/internal/models/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	BotAPI     *tgbotapi.BotAPI
	cfg        *config.Config
	tgUC       bot.Usecase
	userStates map[int64]models.AddQuestionParams
}

func NewTgBot(
	cfg *config.Config,
	usecase bot.Usecase,
	botAPI *tgbotapi.BotAPI,
) *TgBot {
	return &TgBot{
		cfg:        cfg,
		BotAPI:     botAPI,
		tgUC:       usecase,
		userStates: make(map[int64]models.AddQuestionParams),
	}
}

func (t *TgBot) Run() error {
	log.Printf("Authorized on account %s", t.BotAPI.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := t.BotAPI.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {

			switch update.Message.Command() {
			case startCommand:
				if err := t.handleStartCommand(
					update.Message.Chat.ID,
					models.UserActivation{ChatID: update.Message.Chat.ID, TgName: update.Message.Chat.UserName},
					update.Message.Text,
				); err != nil {
					log.Printf("bot.TgBot.handleStartCommand: %s", err.Error())
					if strings.Contains(err.Error(), "Wrong number of rows affected") {
						t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errUUIDAlreadyExists)
						continue
					}
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errUserActivation)
					continue
				}
				continue
			case getUUIDCommand:
				if err := t.handleGetUUIDCommand(
					update.Message.Chat.ID,
					models.GetUUID{ChatID: update.Message.Chat.ID, TgName: update.Message.Chat.UserName}); err != nil {
					log.Printf("bot.TgBot.handleGetUUIDCommand: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
				continue
			case askMeCommend:
				if err := t.handleAskMeCommand(
					update.Message.Chat.ID,
					models.AskMeParams{ChatID: update.Message.Chat.ID}); err != nil {
					log.Printf("bot.TgBot.handleAskMeCommand: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
				continue
			case addQuestion:
				if err := t.handleAddQuestionCommand(
					update.Message.Chat.ID); err != nil {
					log.Printf("bot.TgBot.handleAddQuestionCommand: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
				continue
			}

			questionParams, ok := t.userStates[update.Message.Chat.ID]
			if !ok || questionParams.State == idle {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, addQuestionMessage)
				if _, err := t.BotAPI.Send(msg); err != nil {
					log.Println(err)
				}
				continue
			}

			switch questionParams.State {
			case awaitingQuestion:
				if err := t.handleEnteredQuestion(update.Message.Chat.ID, update.Message.Text); err != nil {
					log.Printf("bot.TgBot.handleEnteredQuestion: %s", err.Error())
					t.userStates[update.Message.Chat.ID] = models.AddQuestionParams{State: idle}
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case awaitingAnswer:
				if err := t.handleEnteredAnswer(update.Message.Chat.ID, update.Message.Text); err != nil {
					log.Printf("bot.TgBot.handleEnteredAnswer: %s", err.Error())
					t.userStates[update.Message.Chat.ID] = models.AddQuestionParams{State: idle}
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			}

		}

		if update.CallbackQuery != nil {
			callbackData, err := t.extractCallbackData(update.CallbackQuery.Data)
			if err != nil {
				log.Println(err)
				continue
			}

			chatID := update.CallbackQuery.From.ID
			messageID := update.CallbackQuery.Message.MessageID
			switch callbackData[0] {
			case backendCallbackData:
				if err := t.handleBackendCallbackData(chatID, messageID); err != nil {
					log.Printf("bot.TgBot.handleBackendCallbackData: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case frontednCallbackData:
				if err := t.handleFrontendCallbackData(chatID, messageID); err != nil {
					log.Printf("bot.TgBot.handleFrontendCallbackData: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case getAnswerCallbackData:
				if err := t.handleGetAnswerCallbackData(chatID, callbackData[1], messageID); err != nil {
					log.Printf("bot.TgBot.handleGetAnswerCallbackData: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case GoCallbackDataAskMe:
				if err := t.handleSubdirectionsCallback(chatID, GoCallbackDataAskMe[:1]); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case ComputerSinceCallbackDataAskMe:
				if err := t.handleSubdirectionsCallback(chatID, ComputerSinceCallbackDataAskMe[:1]); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case NetworkCallbackDataAskMe:
				if err := t.handleSubdirectionsCallback(chatID, NetworkCallbackDataAskMe[:1]); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case DBCallbackDataAskMe:
				if err := t.handleSubdirectionsCallback(chatID, DBCallbackDataAskMe[:1]); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case AlgorithmsCallbackDataAskMe:
				if err := t.handleSubdirectionsCallback(chatID, AlgorithmsCallbackDataAskMe[:1]); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case ArchitectureCallbackDataAskMe:
				if err := t.handleSubdirectionsCallback(chatID, ArchitectureCallbackDataAskMe[:1]); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case GeneralCallbackDataAskMe:
				if err := t.handleSubdirectionsCallback(chatID, GeneralCallbackDataAskMe[:1]); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			}
		}
	}
	return nil
}

func (t *TgBot) extractCallbackData(callbackData string) (result []string, err error) {
	splitedCallbackData := strings.Split(callbackData, " ")
	switch len(splitedCallbackData) {
	case 1:
		return splitedCallbackData, nil
	case 2:
		return splitedCallbackData, nil
	default:
		err = fmt.Errorf("Wrong len of callback data")
		return
	}
}

func (t *TgBot) hideKeyboard(chatID int64, messageID int) (err error) {
	edit := tgbotapi.NewEditMessageReplyMarkup(
		chatID,
		messageID,
		tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
		},
	)

	if _, err = t.BotAPI.Send(edit); err != nil {
		return
	}

	return
}

func (t *TgBot) createDirectionsKeyboard() (keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backenddButton, backendCallbackData),
			tgbotapi.NewInlineKeyboardButtonData(frontendButton, frontednCallbackData),
		),
	)

	return
}

func (t *TgBot) sendErrorMessage(ctx context.Context, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := t.BotAPI.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (t *TgBot) createMainMenuKeyboard() (keyboard tgbotapi.ReplyKeyboardMarkup) {

	keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(getUUIDButton),
			tgbotapi.NewKeyboardButton(askMeButton),
			tgbotapi.NewKeyboardButton(addQuestionButton),
		),
	)

	keyboard.OneTimeKeyboard = false // Hide keyboard after one use
	keyboard.ResizeKeyboard = true   // Resizes keyboard depending on the user's device

	return
}
