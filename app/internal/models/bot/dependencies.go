package bot

import (
	"boost-my-skills-bot/pkg/logger"
	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Dependencies struct {
	PgDB     *sqlx.DB
	Logger   *logger.Logger
	RabbitMQ RabbitMQ
}

type RabbitMQ struct {
	Conn   *amqp.Connection
	Chann  *amqp.Channel
	Queues Queues
}

type Queues struct {
	UserActivationQueue amqp.Queue
}
