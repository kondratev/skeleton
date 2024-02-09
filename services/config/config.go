package config

import (
	env "github.com/caarlos0/env/v7"
	di "github.com/samber/do/v2"

	"github.com/kondratev/skeleton/services/logger"
)

type Service struct {
	log  *logger.Service
	Db   DbConfig   `envPrefix:"DB_"`
	HTTP HTTPConfig `envPrefix:"HTTP_"`
}

func Load(cfg *Service) error {
	parseOpts := env.Options{
		Prefix: "SKEL_",
	}

	err := env.Parse(cfg, parseOpts)
	return err
}

func New(i di.Injector) (*Service, error) {
	log := di.MustInvoke[*logger.Service](i)
	cfg := Service{}
	cfg.log = log
	err := Load(&cfg)
	if err != nil {
		log.Fatal().Err(err)
	}
	return &cfg, nil
}

func (s *Service) Start() {
	s.log.Info().Msg("Config service started")
}

func (s *Service) Reload() error {
	err := Load(s)
	if err != nil {
		s.log.Fatal().Err(err)
		return err
	}
	return nil
}

func (s *Service) Shutdown() error {
	s.log.Info().Msg("Config service stopped")
	return nil
}
