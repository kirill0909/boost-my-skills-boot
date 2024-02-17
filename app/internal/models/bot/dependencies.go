package bot

import (
	"github.com/jmoiron/sqlx"
	"github.com/kirill0909/logger"
)

type Dependencies struct {
	PgDB   *sqlx.DB
	Logger *logger.Logger
}
