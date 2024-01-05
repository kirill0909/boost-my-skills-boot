package bot

import (
	"github.com/jmoiron/sqlx"
	"github.com/kirill0909/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Dependencies struct {
	PgDB     *sqlx.DB
	Logger   *logger.Logger
	RabbitMQ RabbitMQ
}

type RabbitMQ struct {
	Producer Producer
}

type Producer struct {
	Conn   *amqp.Connection
	Chann  *amqp.Channel
	Queues Queues
}

type Queues struct {
	UserActivationQueue amqp.Queue
}
