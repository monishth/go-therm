package core

import (
	"context"
	"log"
	"time"
)

func (a *App) RunControllerLoop() context.CancelFunc {
	log.Println("Controllers initialised")
	ticker := time.NewTicker(time.Minute * 1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("Context cancelled, stopping engine")
				return
			case <-ticker.C:
				a.controllerTick()
			}
		}
	}()

	return cancel
}

func (a *App) controllerTick() {
	for zoneID, controller := range a.Controllers {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

		measuredTemp, err := a.TimeSeriesDataStore.FetchZoneTemp(zoneID, ctx)
		if err != nil {
			cancel()
			continue
		}

		log.Printf("Got measurement for zoneID: %d %v", zoneID, measuredTemp)
		output := controller.Calculate(measuredTemp)
		a.TimeSeriesDataStore.WriteTemperatureTarget(zoneID, controller.setpoint, controller.setpoint-measuredTemp, output)
		log.Printf("Calculated output: %f", output)

		// Kill context
		cancel()
	}
}
