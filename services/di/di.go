package di

import (
	di "github.com/samber/do"
)

type Service interface {
	Get() *di.Injector
}

func New(i *di.Injector) (Service, error) {
	s := ServiceImpl{
		injector: i,
	}
	return &s, nil
}

type ServiceImpl struct {
	injector *di.Injector
}

func (s *ServiceImpl) Start() {
}

func (s *ServiceImpl) Get() *di.Injector {
	return s.injector
}

func (s *ServiceImpl) Shutdown() error {
	println("Injector servive stopped")
	return nil
}
