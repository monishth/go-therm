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
	Targets             map[int]float64
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

// These should move

func (a *App) SetTarget(zoneID int, target float64) {
	a.Targets[zoneID] = target
	a.TimeSeriesDataStore.WriteTarget(zoneID, target)
}
