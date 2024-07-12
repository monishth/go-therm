package engine

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/monishth/go-therm/internal/db"
	"github.com/monishth/go-therm/internal/models"
	"github.com/monishth/go-therm/internal/mqtt"
)

// All this does is store temp measurements atm
func RunEngine(wg *sync.WaitGroup) {
	defer wg.Done()

	config := db.LoadConfig()
	mqttClient := mqtt.StartMQTTClient()

	for _, thermostat := range config.Thermostats {
		handler := func(state models.SensorState) {
			writeAPI := db.GetInfluxClient()
			db.WriteTemperature(writeAPI, thermostat.ZoneID, thermostat.Name, state.DS18B20.Temperature)
		}

		mqtt.SubscribeToTopic(mqttClient, thermostat.Topic, mqtt.WrapHandler(handler))
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	db.CloseInfluxClient()
	mqtt.ShutdownMQTTClient(mqttClient)
	log.Println("Graceful shutdown completed")
}
