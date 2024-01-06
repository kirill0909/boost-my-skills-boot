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

	if err = queueDeclare(ch, cfg); err != nil {
		return models.Producer{}, err
	}

	return models.Producer{
		Conn:  conn,
		Chann: ch}, nil
}

func queueDeclare(ch *amqp.Channel, cfg *config.Config) error {
	queueNames := []string{
		cfg.RabbitMQ.QueueNames.UserActivationQueue,
		cfg.RabbitMQ.QueueNames.GetUpdatedButtonsQueue}

	for _, name := range queueNames {
		if _, err := ch.QueueDeclare(
			name,  // name
			false, // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		); err != nil {
			err = errors.Wrap(err, "rabbit.queueDeclare. Failed to declare queue")
			return err
		}
	}

	return nil

}
