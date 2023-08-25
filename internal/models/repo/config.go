package repo

import (
	"github.com/jmoiron/sqlx"

	"github.com/kondratev/skeleton/services/db"
	"github.com/kondratev/skeleton/services/logger"
)

type Config struct {
	Log *logger.Service
	Db  db.Service
	SQL *sqlx.DB
}
