package httpmux

import (
	"context"
	"math/rand"
	"net"
	"time"

	ginlogger "github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid"
	"github.com/rs/zerolog"
	di "github.com/samber/do/v2"

	"github.com/kondratev/skeleton/internal/consts"
	"github.com/kondratev/skeleton/services/config"
	"github.com/kondratev/skeleton/services/logger"
)

type Service interface {
	Ping() error
	Start() error
	Shutdown() error
	Server() *gin.Engine
}

func New(i di.Injector) (Service, error) {
	cfg := di.MustInvoke[*config.Service](i)
	log := di.MustInvoke[*logger.Service](i)

	mux := gin.New()

	mux.Use(gin.Recovery())

	mux.Use(
		requestid.New(
			requestid.WithGenerator(func() string {
				return ulid.MustNew(
					ulid.Timestamp(time.Now()),
					rand.New(rand.NewSource(time.Now().UnixNano()))).
					String()
			}),
			requestid.WithCustomHeaderStrKey(consts.ULIDHeader),
		),
	)

	mux.Use(ginlogger.SetLogger(
		ginlogger.WithLogger(
			func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
				return l.Output(gin.DefaultWriter).With().Str("request_id", requestid.Get(c)).Logger()
			})))

	s := ServiceImpl{
		mux: mux,
		cfg: &cfg.HTTP,
		log: log,
	}

	return &s, nil
}

type ServiceImpl struct {
	mux      *gin.Engine
	listener net.Listener
	cfg      *config.HTTPConfig
	log      *logger.Service
	ctx      context.Context
	cancel   context.CancelFunc
}

func (s *ServiceImpl) Start() error {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	lc := net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", s.cfg.ListenAddr)
	if err != nil || listener == nil {
		cancel()
		return err
	}
	s.listener = listener
	s.cancel = cancel
	s.ctx = ctx

	s.log.Info().Msg("HTTP service started")

	go func() {
		err = s.mux.RunListener(s.listener)
		if err != nil {
			s.log.Err(err).Msg("can't start server httpmux on listener")
		}
		s.cancel()
	}()
	return nil
}

func (s *ServiceImpl) Shutdown() error {
	if err := s.listener.Close(); err != nil {
		s.log.Warn().Err(err)
	}
	s.cancel()
	s.log.Info().Msg("HTTP service stopped")
	return nil
}

func (s *ServiceImpl) HealthCheck() error {
	return s.Ping()
}

func (s *ServiceImpl) Ping() error {
	return nil
}

func (s *ServiceImpl) Server() *gin.Engine {
	return s.mux
}

func (s *ServiceImpl) Context() context.Context {
	return s.ctx
}
