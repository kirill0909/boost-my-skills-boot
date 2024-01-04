package rabbit

import (
	"boost-my-skills-bot/config"
	models "boost-my-skills-bot/internal/models/bot"
	"fmt"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbit(cfg *config.Config) (models.RabbitMQ, error) {
	url := fmt.Sprintf("%s:%s/", cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		err = errors.Wrap(err, "rabbit.InitRabbit. failed to connect to rabbitMQ")
		return models.RabbitMQ{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		err = errors.Wrap(err, "rabbit.InitRabbit. failed to open channel")
		return models.RabbitMQ{}, err
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
		err = errors.Wrap(err, "rabbit.InitRabbit. Failed to declare queue ")
		return models.RabbitMQ{}, err
	}

	return models.RabbitMQ{
			Conn:  conn,
			Chann: ch,
			Queues: models.Queues{
				UserActivationQueue: userActivationQueue,
			},
		},
		nil
}
