package engine

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/monishth/go-therm/internal/db"
	"github.com/monishth/go-therm/internal/models"
	"github.com/monishth/go-therm/internal/mqtt"
	"github.com/monishth/go-therm/pkg/utils"
)

// All this does is store temp measurements atm
func RunEngine(wg *sync.WaitGroup) {
	defer wg.Done()

	config := db.LoadConfig()
	mqttClient := mqtt.StartMQTTClient()

	subscribeThermostats(config, mqttClient)

	subscribeValves(config, mqttClient)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	db.CloseInfluxClient()
	mqtt.ShutdownMQTTClient(mqttClient)
	log.Println("Graceful shutdown completed")
}

func subscribeThermostats(config db.Config, mqttClient MQTT.Client) {
	for _, thermostat := range config.Thermostats {
		handler := func(state models.SensorState) {
			writeAPI := db.GetInfluxClient()
			db.WriteTemperature(writeAPI, thermostat.ZoneID, thermostat.ID, state.DS18B20.Temperature)
		}

		mqtt.SubscribeToTopic(mqttClient, thermostat.Topic, mqtt.WrapHandler(handler))
	}
}

func subscribeValves(config db.Config, mqttClient MQTT.Client) {
	topicToHandlers := make(map[string][]func(map[string]interface{}))
	for _, valve := range config.Valves {
		handler := func(state map[string]interface{}) {
			writeAPI := db.GetInfluxClient()
			db.WriteValveState(writeAPI, valve.ZoneID, valve.ID, utils.ConvertOnOffToInt(state[valve.RelayName].(string)))
		}
		topicToHandlers[valve.StateTopic] = append(topicToHandlers[valve.StateTopic], handler)
	}

	for key, value := range topicToHandlers {
		log.Printf("Registering Topic: %s, Handlers: %d", key, len(value))
		mqtt.SubscribeToTopic(mqttClient, key, mqtt.WrapHandler(func(state map[string]interface{}) {
			for _, handler := range value {
				handler(state)
			}
		}))
	}
}
