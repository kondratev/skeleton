package logger

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	di "github.com/samber/do/v2"
)

type Service struct {
	zerolog.Logger
}

func New(di.Injector) (*Service, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	serv := Service{log}
	log.Info().Msg("Logger service started")
	return &serv, nil
}

func (s *Service) Shutdown() error {
	println(`{"level":"info","time":` + strconv.FormatInt(time.Now().Unix(), 10) + `,"message":"Logger service stopped"}`)
	return nil
}

func (s *Service) IfErr(err error, msg string) error {
	if err != nil {
		s.Logger.Err(err).Msg(msg)
	}
	return err
}

func (s *Service) IfErrLog(err error, msg string) {
	if err != nil {
		s.Logger.Err(err).Msg(msg)
	}
}

func (s *Service) IfErrf(err error, format string, v ...interface{}) error {
	if err != nil {
		s.Logger.Err(err).Msgf(format, v...)
	}
	return err
}

func (s *Service) IfErrFatal(err error, msg string) {
	if err != nil {
		s.Logger.Fatal().Err(err).Msg(msg)
	}
}

func (s *Service) IfErrFatalf(err error, msg string, v ...interface{}) {
	if err != nil {
		s.Logger.Fatal().Err(err).Msgf(msg, v...)
	}
}
