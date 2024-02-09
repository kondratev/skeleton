package closer

import (
	"os"

	di "github.com/samber/do/v2"
)

type Service struct{}

func New(di.Injector) (*Service, error) {
	serv := Service{}
	return &serv, nil
}

func (s *Service) Start() {
}

func (s *Service) Shutdown() error {
	os.Exit(0)
	return nil
}
