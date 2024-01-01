package bot

import (
	"boost-my-skills-bot/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	PgDB   *sqlx.DB
	Logger *logger.Logger
}
