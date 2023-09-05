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
					update.Message.Chat.ID, update.Message.Text,
					questionParams.SubdirectionID, questionParams.SubSubdirectionID); err != nil {
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
			// Ask me sub callbacks
			case callbackDataAskMe[0]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[0], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataAskMe[1]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataAskMe[2]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[2], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataAskMe[3]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[3], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataAskMe[4]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[4], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataAskMe[5]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[5], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataAskMe[6]:
				if err := t.handleSubdirectionsCallbackAskMe(chatID, callbackDataAskMe[6], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
				// Ask me sub sub direction
			case callbackDataSubSubdirectionAskMe[0]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[0], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[1]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[2]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[2], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[3]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[3], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[4]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[4], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[5]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[5], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[6]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[6], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[7]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[7], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[8]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[8], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[9]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[9], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[10]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[10], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[11]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[11], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[12]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[12], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[13]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[13], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAskMe[14]:
				if err := t.handleCallbackDataSubSubdirectionAskMe(chatID, callbackDataSubSubdirectionAskMe[14], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallback: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
				// Add SUB directions question callbacks
			case callbackDataSubdirectionAddQuestion[0]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[0], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubdirectionAddQuestion[1]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubdirectionAddQuestion[2]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[2], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubdirectionAddQuestion[3]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[3], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubdirectionAddQuestion[4]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[4], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubdirectionAddQuestion[5]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[5], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubdirectionAddQuestion[6]:
				if err := t.handleSubdirectionsCallbackAddQuestion(chatID, callbackDataSubdirectionAddQuestion[6], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
				// Add SUB SUB directions question callbacks
			case callbackDataSubSubdirectionAddQuestion[0]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[0], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[1]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[1], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[2]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[2], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[3]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[3], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[4]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[4], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[5]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[5], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[6]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[6], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[7]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[7], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[8]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[8], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[9]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[9], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[10]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[10], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[11]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[11], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[12]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[12], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[13]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[13], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[14]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[14], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[15]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[15], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[16]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[16], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[17]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[17], messageID); err != nil {
					log.Printf("bot.TgBot.handleSubSubdirectionsCallbackAddQuestion: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errInternalServerError)
					continue
				}
			case callbackDataSubSubdirectionAddQuestion[18]:
				if err := t.handleSubSubdirectionsCallbackAddQuestion(
					chatID, callbackDataSubSubdirectionAddQuestion[18], messageID); err != nil {
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
	n := len(splitedCallbackData)
	switch {
	case n == 1 || n == 2:
		return splitedCallbackData, nil
	default:
		err = fmt.Errorf("Wrong lenngth(%d) of callback data", n)
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
