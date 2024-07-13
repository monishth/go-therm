package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/monishth/go-therm/internal/core"
	"github.com/monishth/go-therm/internal/frontend"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	app := core.CreateApp()
	app.Listen()
	app.InitialiseControllers()

	app.RunControllerLoop(ctx)
	go frontend.StartFrontend(&app, ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	cancel()
	app.Shutdown()
}
