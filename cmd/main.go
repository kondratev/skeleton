package main

import (
	"github.com/samber/do/v2"
)

func main() {
	injector := do.New()
	registerProviders(injector)
	startServices(injector)
	injector.ShutdownOnSignals()
}
