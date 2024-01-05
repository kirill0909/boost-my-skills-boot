package usecase

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kirill0909/logger"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
	"time"
)

type BotUC struct {
	cfg                  *config.Config
	pgRepo               bot.PgRepository
	rabbitMQ             models.RabbitMQ
	BotAPI               *tgbotapi.BotAPI
	log                  *logger.Logger
	lastKeyboardChecking int64
}

func NewBotUC(
	cfg *config.Config,
	pgRepo bot.PgRepository,
	rabbitMQ models.RabbitMQ,
	botAPI *tgbotapi.BotAPI,
	log *logger.Logger,
) bot.Usecase {
	return &BotUC{
		cfg:                  cfg,
		pgRepo:               pgRepo,
		rabbitMQ:             rabbitMQ,
		BotAPI:               botAPI,
		log:                  log,
		lastKeyboardChecking: time.Now().Unix(),
	}
}

func (u *BotUC) HandleStartCommand(ctx context.Context, params models.HandleStartCommandParams) error {
	splitedText := strings.Split(params.Text, " ")

	if len(splitedText) != 2 {
		return fmt.Errorf("TgBot.handleStartCommand. wrong len of splited text: %d != 2. params(%+v)", len(splitedText), params)
	}

	uuid := splitedText[1]
	setStatusActiveParams := models.UserActivationParams{
		TgName: params.TgName,
		ChatID: params.ChatID,
		UUID:   uuid}

	setStatusActiveParamsBytes, err := json.Marshal(setStatusActiveParams)
	if err != nil {
		err = errors.Wrapf(err, "BotUC.HandleStartCommand.Marshal. params(%+v)", setStatusActiveParams)
	}

	if err := u.writeToBroker(u.rabbitMQ.Queues.UserActivationQueue.Name, setStatusActiveParamsBytes); err != nil {
		return err
	}

	var isAdmin bool
	msg := tgbotapi.NewMessage(params.ChatID, "your account has been successfully activated")
	if params.ChatID == u.cfg.AdminChatID {
		isAdmin = true
		msg.ReplyMarkup = u.createMainMenuKeyboard(isAdmin)
	} else {
		isAdmin = false
		msg.ReplyMarkup = u.createMainMenuKeyboard(isAdmin)
	}

	if _, err := u.BotAPI.Send(msg); err != nil {
		return errors.Wrapf(err, "BotUC.HandleStartCommand.Send")

	}

	return nil
}

func (u *BotUC) writeToBroker(queue string, message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.rabbitMQ.Chann.PublishWithContext(ctx,
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message})
	if err != nil {
		return errors.Wrapf(err, "unable to wirte message(%s) to queue(%s)", message, queue)
	}

	u.log.Infof("sent message(%s) to queue(%s)", message, queue)

	return nil
}
