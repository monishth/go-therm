package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/monishth/go-therm/internal/core"
)

func main() {
	app := core.CreateApp()
	app.Listen()
	app.InitialiseControllers()

	tickerCancel := app.RunControllerLoop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	tickerCancel()
	app.Shutdown()
	log.Println("Engine shutdown")
}
