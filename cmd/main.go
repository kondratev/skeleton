package main

import (
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/samber/do"
)

func main() {
	injector := do.New()
	registerProviders(injector)
	startServices(injector)
	log.Fatal().Err(injector.ShutdownOnSignals(syscall.SIGINT))
}
