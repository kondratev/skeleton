package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	di "github.com/samber/do/v2"

	"github.com/kondratev/skeleton/services/config"
	"github.com/kondratev/skeleton/services/httpmux"
	"github.com/kondratev/skeleton/services/logger"
	"github.com/kondratev/skeleton/services/repository"
)

const LivenessRoute = "/ping"

type Service interface {
	Start()
	Shutdown() error
}

type ServiceImpl struct {
	log  *logger.Service
	cfg  *config.Service
	http *gin.Engine
	repo repository.Service
}

func New(i di.Injector) (Service, error) {
	loggerService := di.MustInvoke[*logger.Service](i)
	cfg := di.MustInvoke[*config.Service](i)
	httpService := di.MustInvoke[httpmux.Service](i)
	repo := di.MustInvoke[repository.Service](i)

	return &ServiceImpl{
		log:  loggerService,
		cfg:  cfg,
		http: httpService.Server(),
		repo: repo,
	}, nil
}

func (s *ServiceImpl) Start() {
	s.http.GET(LivenessRoute, s.Ping())
	s.http.GET("/metrics", gin.WrapH(promhttp.Handler()))

	s.log.Info().Msg("HttpHandler service started")
}

func (s *ServiceImpl) Shutdown() error {
	s.log.Info().Msg("HttpHandler service stopped")
	return nil
}

func (s *ServiceImpl) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		type pingResp struct {
			Message string
		}
		c.JSON(http.StatusOK, &pingResp{Message: "pong"})
		s.log.Info().Msg("ping request")
	}
}
