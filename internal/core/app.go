package core

import (
	"log"
	"sync"

	"github.com/monishth/go-therm/internal/db"
	"github.com/monishth/go-therm/internal/messaging"
)

type App struct {
	Entities            db.Entities
	MessageClient       messaging.MessageClient
	TimeSeriesDataStore db.TimeSeriesDataStore
	Controllers         map[int]*PIDController
	targets             map[int]float64
	targetsMutex        sync.RWMutex
}

func (a *App) GetTargets() map[int]float64 {
	a.targetsMutex.RLock()
	defer a.targetsMutex.RUnlock()
	return a.targets
}

// All this does is store temp measurements atm
func CreateApp() App {
	entities := db.LoadEntities()
	mqttClient := messaging.StartMQTTClient()
	dbClient := db.CreateInfluxClient()

	return App{
		Entities:            entities,
		MessageClient:       &mqttClient,
		TimeSeriesDataStore: dbClient,
	}
}

func (a *App) Listen() {
	// Start Listeners
	a.subscribeThermostats()
	a.subscribeValves()
}

func (a *App) Shutdown() {
	a.TimeSeriesDataStore.Close()
	a.MessageClient.Close()
	log.Println("Engine shutdown")
}

// This should probably move
func (a *App) SetTarget(zoneID int, target float64) {
	a.targetsMutex.Lock()
	defer a.targetsMutex.Unlock()
	a.targets[zoneID] = target
	a.TimeSeriesDataStore.WriteTarget(zoneID, target)
}
