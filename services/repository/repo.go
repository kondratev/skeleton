package repository

import (
	"github.com/jmoiron/sqlx"
	di "github.com/samber/do/v2"

	"github.com/kondratev/skeleton/internal/models/repo"
	"github.com/kondratev/skeleton/services/db"
	"github.com/kondratev/skeleton/services/logger"
)

type Service interface {
	Start() error
	HealthCheck() error
	CreateTransaction() *sqlx.Tx
}

type ServiceImpl struct {
	log *logger.Service
	db  db.Service
	sql *sqlx.DB
}

func New(i di.Injector) (Service, error) {
	loggerService := di.MustInvoke[*logger.Service](i)
	dbService := di.MustInvoke[db.Service](i)

	return &ServiceImpl{
		log: loggerService,
		db:  dbService,
	}, nil
}

func (s *ServiceImpl) Start() error {
	s.sql = s.db.Db()

	s.log.Info().Msg("Repository service started")

	return nil
}

func (s *ServiceImpl) config() *repo.Config {
	return &repo.Config{
		Log: s.log,
		Db:  s.db,
		SQL: s.sql,
	}
}

func (s *ServiceImpl) Shutdown() error {
	s.log.Info().Msg("Repository service stopped")
	return nil
}

func (s *ServiceImpl) HealthCheck() error {
	return nil
}

func (s *ServiceImpl) CreateTransaction() *sqlx.Tx {
	return s.sql.MustBegin()
}
