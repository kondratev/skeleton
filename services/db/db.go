package db

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx postgresql driver
	"github.com/jmoiron/sqlx"
	di "github.com/samber/do"

	"github.com/kondratev/skeleton/services/config"
	"github.com/kondratev/skeleton/services/logger"
)

type Service interface {
	HealthCheck() error
	Start() error
	Db() *sqlx.DB
	Log() *logger.Service
}

func (s *ServiceImpl) Log() *logger.Service {
	return s.log
}

const baseDriverName = "postgres"

type ServiceImpl struct {
	log    *logger.Service
	db     *sqlx.DB
	ctx    *context.Context
	dsname *string
	cfg    *config.Service
}

func New(i *di.Injector) (Service, error) {
	log := di.MustInvoke[*logger.Service](i)
	cfg := di.MustInvoke[*config.Service](i)

	ctx := context.Background()

	return &ServiceImpl{
		log: log,
		ctx: &ctx,
		cfg: cfg,
	}, nil
}

type User struct {
	ID        int    `db:"id" goqu:"skipinsert"`
	FirstName string `db:"tex"`
	Age       int    `db:"myint"`
}

func (s *ServiceImpl) Start() error {
	var dsn string

	if len(s.cfg.Db.DsName) == 0 {
		return errors.New("datasource name not set")
	}
	if s.cfg.Db.Driver == "pgx" || s.cfg.Db.Driver == baseDriverName {
		dsn = fmt.Sprintf("%s://%s:%s@%s:%s/%s",
			baseDriverName,
			s.cfg.Db.User,
			s.cfg.Db.Password,
			s.cfg.Db.Host,
			s.cfg.Db.Port,
			s.cfg.Db.Database)
		s.cfg.Db.Driver = "pgx"
		s.dsname = &s.cfg.Db.DsName
	} else {
		s.log.Fatal().Err(errors.New("env not set property")).Msgf("database driver not set, ds:%s", s.cfg.Db.DsName)
	}
	db, err := sqlx.Connect(s.cfg.Db.Driver, dsn)
	if err := s.log.IfErrf(err, "can't open database, ds:%s", *s.dsname); err != nil {
		return err
	}
	if err := s.log.IfErrf(db.Ping(), "can't connect to database, ds:%s", *s.dsname); err != nil {
		return err
	}
	s.db = db

	s.log.Info().Msgf("Db service started, ds:%s", *s.dsname)
	return nil
}

func (s *ServiceImpl) Shutdown() error {
	s.log.IfErrLog(s.db.Close(), "close db")
	s.log.Info().Msgf("Db service stopped, ds:%s", *s.dsname)
	return nil
}

func (s *ServiceImpl) HealthCheck() error {
	return s.db.Ping()
}

func (s *ServiceImpl) Db() *sqlx.DB {
	return s.db
}
