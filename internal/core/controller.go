package core

import (
	"context"
	"log"
	"time"
)

func (a *App) InitialiseControllers() {
	if a.Controllers == nil {
		controllers := make(map[int]*PIDController)
		targets := make(map[int]float64)

		for _, zone := range a.Config.Zones {
			log.Printf("Creating PID Controller for Zone: %d", zone.ID)
			controllers[zone.ID] = NewPIDController(5, 0.5, 1, 21)
			targets[zone.ID] = 21
		}
		a.Controllers = controllers
		a.Targets = targets
	}
}

func (a *App) RunControllerLoop(ctx context.Context) {
	log.Println("Controllers initialised")
	ticker := time.NewTicker(time.Minute * 1)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("Context cancelled, stopping engine")
				return
			case <-ticker.C:
				a.controllerTick(ctx)
			}
		}
	}()
}

func (a *App) controllerTick(ctx context.Context) {
	for zoneID, controller := range a.Controllers {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*20)

		measuredTemp, err := a.TimeSeriesDataStore.FetchZoneTemp(zoneID, ctxWithTimeout)
		if err != nil {
			cancel()
			continue
		}

		target := a.Targets[zoneID]

		controller.SetSetpoint(target)

		log.Printf("Got measurement for zoneID: %d %v %v", zoneID, measuredTemp, target)
		output := controller.Calculate(measuredTemp)
		a.TimeSeriesDataStore.WriteTemperatureData(zoneID, controller.setpoint, controller.setpoint-measuredTemp, output)
		log.Printf("Calculated output: %f", output)

		// Kill context
		cancel()
	}
}

const (
	HeatLossRate  = 0.1  // Rate at which the room loses heat (degrees per minute)
	HeatGainRate  = 0.3  // Rate at which the room gains heat (degrees per minute)
	HysteresisOn  = 10.0 // Hysteresis for turning valve on (percentage)
	HysteresisOff = 5.0  // Hysteresis for turning valve off (percentage)
)

// ActivationFunction determines whether the valve should be on or off
func ActivationFunction(pidOutput float64, currentState bool, currentTemp, setpoint float64) bool {
	thresholdOn := HeatLossRate * 100 // Adjust based on your specific characteristics
	thresholdOff := -HeatGainRate * 100

	// Turning on the valve
	if pidOutput > thresholdOn+HysteresisOn {
		return true
	}

	// Turning off the valve
	if pidOutput < thresholdOff-HysteresisOff {
		return false
	}

	// Hysteresis to prevent rapid switching
	if currentState && pidOutput < -HysteresisOff {
		return false
	}
	if !currentState && pidOutput > HysteresisOn {
		return true
	}

	return currentState
}
