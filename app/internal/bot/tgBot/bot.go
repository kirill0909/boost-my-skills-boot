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
			if !ok || questionParams.State == idle ||
				questionParams.State == awaitingSubdirection ||
				questionParams.State == awaitingSubSubdirection {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, addQuestionMessage)
				if _, err := t.BotAPI.Send(msg); err != nil {
					log.Println(err)
				}
				continue
			}

			switch questionParams.State {
			case awaitingQuestion:
				if err := t.handleEnteredQuestion(
					update.Message.Chat.ID, update.Message.Text, questionParams.SubdirectionID); err != nil {
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
			// Ask me callbacks
			case callbackDataAskMe[0]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[0][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataAskMe[1]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[1][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataAskMe[2]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[2][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataAskMe[3]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[3][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataAskMe[4]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[4][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataAskMe[5]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[5][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataAskMe[6]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[6][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
				// Add SUB directions question callbacks
			case callbackDataSubdirectionAddQuestion[0]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[0][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubdirectionAddQuestion[1]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[1][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubdirectionAddQuestion[2]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[2][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubdirectionAddQuestion[3]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[3][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubdirectionAddQuestion[4]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[4][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubdirectionAddQuestion[5]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[5][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubdirectionAddQuestion[6]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[6][:1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
				// Add SUB SUB directions question callbacks
			case callbackDataSubSubdirectionAddQuestion[0]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[0][:1], callbackDataSubSubdirectionAddQuestion[0][2:3], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[1]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[1][:1], callbackDataSubSubdirectionAddQuestion[1][2:3], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
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

func (t *TgBot) sendErrorMessage(ctx context.Context, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := t.BotAPI.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
