package tgbot

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	"log"
	"strconv"
	"strings"

	models "boost-my-skills-bot/internal/models/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	BotAPI          *tgbotapi.BotAPI
	cfg             *config.Config
	tgUC            bot.Usecase
	stateUsers      map[int64]models.AddInfoParams
	stateDirections *models.DirectionsData
}

func NewTgBot(
	cfg *config.Config,
	usecase bot.Usecase,
	botAPI *tgbotapi.BotAPI,
	stateUsers map[int64]models.AddInfoParams,
	stateDirections *models.DirectionsData,
) *TgBot {
	return &TgBot{
		cfg:             cfg,
		BotAPI:          botAPI,
		tgUC:            usecase,
		stateUsers:      stateUsers,
		stateDirections: stateDirections,
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
				params := models.UserActivation{
					ChatID: update.Message.Chat.ID, TgName: update.Message.Chat.UserName, Text: update.Message.Text}
				if err := t.handleStartCommand(params); err != nil {
					log.Printf("bot.TgBot.handleStartCommand: %s", err.Error())
					if strings.Contains(err.Error(), "Wrong number of rows affected") {
						t.sendErrorMessage(update.Message.Chat.ID, errUUIDAlreadyExists)
						continue
					}
					t.sendErrorMessage(update.Message.Chat.ID, errUserActivation)
					continue
				}
				continue
			case getUUIDCommand:
				if err := t.handleGetUUIDCommand(
					update.Message.Chat.ID,
					models.GetUUID{ChatID: update.Message.Chat.ID, TgName: update.Message.Chat.UserName}); err != nil {
					log.Printf("bot.TgBot.handleGetUUIDCommand: %s", err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
				continue
			case askMeCommend:
				if err := t.handleAskMeCommand(update.Message.Chat.ID, models.AskMeParams{ChatID: update.Message.Chat.ID}); err != nil {
					log.Printf("bot.TgBot.handleAskMeCommand: %s", err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
				continue
			case addInfoCommand:
				if err := t.handleAddInfoCommand(update.Message.Chat.ID); err != nil {
					log.Printf("bot.TgBot.handleAddInfoCommand: %s", err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
				continue
			case printQuestionsCommand:
				if err := t.handlePrintQuestionsCommand(update.Message.Chat.ID); err != nil {
					log.Printf("bot.TgBot.handlePrintInfoCommand: %s", err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
				continue
			}

			questionParams, ok := t.stateUsers[update.Message.Chat.ID]
			if !ok || questionParams.State == t.cfg.StateMachineStatus.Idle ||
				questionParams.State == t.cfg.StateMachineStatus.AwaitingSubdirection ||
				questionParams.State == t.cfg.StateMachineStatus.AwaitingSubSubdirection {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, addQuestionMessage)
				if _, err := t.BotAPI.Send(msg); err != nil {
					log.Println(err)
				}
				continue
			}

			switch questionParams.State {
			case t.cfg.StateMachineStatus.AwaitingQuestion:
				if err := t.handleEnteredQuestion(
					update.Message.Chat.ID, update.Message.Text,
					questionParams.SubdirectionID, questionParams.SubSubdirectionID); err != nil {
					log.Printf("bot.TgBot.handleEnteredQuestion: %s", err.Error())
					t.stateUsers[update.Message.Chat.ID] = models.AddInfoParams{State: t.cfg.StateMachineStatus.Idle}
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case t.cfg.StateMachineStatus.AwaitingAnswer:
				if err := t.handleEnteredAnswer(update.Message.Chat.ID, update.Message.Text); err != nil {
					log.Printf("bot.TgBot.handleEnteredAnswer: %s", err.Error())
					t.stateUsers[update.Message.Chat.ID] = models.AddInfoParams{State: t.cfg.StateMachineStatus.Idle}
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
			}

		}

		if update.CallbackQuery != nil {
			callbackType, err := t.extractCallbackType(update.CallbackQuery.Data)
			if err != nil {
				log.Println(err)
				continue
			}

			chatID := update.CallbackQuery.From.ID
			messageID := update.CallbackQuery.Message.MessageID

			switch callbackType {
			case t.cfg.CallbackType.Direction:
				if err := t.handleDirectionCallbackData(chatID, messageID, update.CallbackQuery.Data); err != nil {
					log.Printf("bot.TgBot.handleDirectionCallbackData: %s", err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case t.cfg.CallbackType.SubdirectionAddInfo:
				if err := t.handleAddInfoSubdirectionCallbackData(chatID, messageID, update.CallbackQuery.Data); err != nil {
					log.Printf("bot.TgBot.handleAddInfoSubdirectionCallbackData: %s", err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case t.cfg.CallbackType.SubSubdirectionAddInfo:
				if err := t.handleAddInfoSubSubdirectionCallbackData(chatID, messageID, update.CallbackQuery.Data); err != nil {
					log.Printf("bot.TgBot.handleAddInfoSubSubdirectionCallbackData: %s", err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case t.cfg.CallbackType.SubdirectionAskMe:
				if err := t.handleAskMeSubdirectionCallbackData(chatID, messageID, update.CallbackQuery.Data); err != nil {
					log.Printf("bot.TgBot.handleAskMeSubdirectionCallbackData: %s", err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case t.cfg.CallbackType.SubSubdirectionAskMe:
				if err := t.handleAskMeSubSubdirectionCallbackData(chatID, messageID, update.CallbackQuery.Data); err != nil {
					log.Printf("bot.TgBot.handleAskMeSubSubdirectionCallbackData: %s", err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case t.cfg.CallbackType.GetAnAnswer:
				if err := t.handleGetAnAnswerCallbackData(chatID, messageID, update.CallbackQuery.Data); err != nil {
					log.Printf("bot.TgBot.handleGetAnAnswerCallbackData: %s", err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case t.cfg.CallbackType.SubdirectionPrintQuestions:
				if err := t.handlePrintQuestionsSubdirectionCallbackData(chatID, messageID, update.CallbackQuery.Data); err != nil {
					log.Printf("bot.TgBot.handlePrintQuestionsSubdirectionCallbackData: %s", err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case t.cfg.CallbackType.SubSubdirectionPrintQuestions:
				if err := t.handlePrintQuestionsSubSubdirectionCallbackData(chatID, messageID, update.CallbackQuery.Data); err != nil {
					log.Printf("bot.TgBot.handlePrintQuestionsSubSubdirectionCallbackData: %s", err.Error())
					t.sendErrorMessage(update.Message.Chat.ID, errInternalServerError)
					continue
				}
			}

		}
	}
	return nil
}

func (t *TgBot) extractCallbackType(callbackData string) (result int, err error) {
	splitedCallbackData := strings.Split(callbackData, " ")
	lastElemCallbackData := splitedCallbackData[len(splitedCallbackData)-1]

	callbackType, err := strconv.Atoi(lastElemCallbackData)
	if err != nil {
		return
	}

	return callbackType, nil
}

func (t *TgBot) sendErrorMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := t.BotAPI.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
