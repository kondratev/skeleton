package main

import (
	"github.com/samber/do"

	"github.com/kondratev/skeleton/services/closer"
	"github.com/kondratev/skeleton/services/config"
	"github.com/kondratev/skeleton/services/db"
	"github.com/kondratev/skeleton/services/di"
	"github.com/kondratev/skeleton/services/handlers"
	"github.com/kondratev/skeleton/services/httpmux"
	"github.com/kondratev/skeleton/services/logger"
	"github.com/kondratev/skeleton/services/repository"
)

func registerProviders(i *do.Injector) {
	do.Provide(i, closer.New)
	do.Provide(i, logger.New)
	do.Provide(i, config.New)
	do.Provide(i, db.New)
	do.Provide(i, handlers.New)
	do.Provide(i, httpmux.New)
	do.Provide(i, di.New)
	do.Provide(i, repository.New)
}

func startServices(i *do.Injector) {
	clo := do.MustInvoke[*closer.Service](i)
	log := do.MustInvoke[*logger.Service](i)
	cfg := do.MustInvoke[*config.Service](i)
	data := do.MustInvoke[db.Service](i)
	repo := do.MustInvoke[repository.Service](i)
	web := do.MustInvoke[handlers.Service](i)
	mux := do.MustInvoke[httpmux.Service](i)

	clo.Start()
	cfg.Start()

	log.IfErrFatalf(data.Start(), "can't init database source, ds:%s", cfg.Db.DsName)
	log.IfErrFatal(repo.Start(), "can't init repository")
	web.Start() // start handlers
	log.IfErrFatal(mux.Start(), "can't start server http mux")
}
