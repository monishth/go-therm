package core

import (
	"log"

	"github.com/monishth/go-therm/internal/db"
	"github.com/monishth/go-therm/internal/messaging"
)

type App struct {
	Config              db.Config
	MessageClient       messaging.MessageClient
	TimeSeriesDataStore db.TimeSeriesDataStore
	Controllers         map[int]*PIDController
}

// All this does is store temp measurements atm
func CreateApp() App {
	config := db.LoadConfig()
	mqttClient := messaging.StartMQTTClient()
	dbClient := db.CreateInfluxClient()

	return App{
		Config:              config,
		MessageClient:       &mqttClient,
		TimeSeriesDataStore: dbClient,
	}
}

func (a *App) InitialiseControllers() {
	if a.Controllers == nil {
		controllers := make(map[int]*PIDController)

		for _, zone := range a.Config.Zones {
			log.Printf("Creating PID Controller for Zone: %d", zone.ID)
			controllers[zone.ID] = NewPIDController(0.1, 0, 0.1, 21)
		}
		a.Controllers = controllers
	}
}

func (a *App) Listen() {
	// Start Listeners
	a.subscribeThermostats()
	a.subscribeValves()
}

func (e *App) Shutdown() {
	e.TimeSeriesDataStore.Close()
	e.MessageClient.Close()
	log.Println("Engine shutdown")
}
