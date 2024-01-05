package rabbit

import (
	"boost-my-skills-bot/config"
	models "boost-my-skills-bot/internal/models/bot"
	"fmt"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitProducer(cfg *config.Config) (models.Producer, error) {
	url := fmt.Sprintf("%s:%s/", cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		err = errors.Wrap(err, "rabbit.InitRabbitProducer. failed to connect to rabbitMQ")
		return models.Producer{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		err = errors.Wrap(err, "rabbit.InitRabbitProducer. failed to open channel")
		return models.Producer{}, err
	}

	userActivationQueue, err := ch.QueueDeclare(
		cfg.RabbitMQ.Queues.UserActivationQueue, // name
		false,                                   // durable
		false,                                   // delete when unused
		false,                                   // exclusive
		false,                                   // no-wait
		nil,                                     // arguments
	)
	if err != nil {
		err = errors.Wrap(err, "rabbit.InitRabbitProducer. Failed to declare queue ")
		return models.Producer{}, err
	}

	return models.Producer{
		Conn:  conn,
		Chann: ch,
		Queues: models.Queues{
			UserActivationQueue: userActivationQueue}}, nil
}
