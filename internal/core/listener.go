package core

import (
	"log"

	"github.com/monishth/go-therm/internal/models"
	"github.com/monishth/go-therm/pkg/utils"
)

func (a *App) subscribeThermostats() {
	for _, thermostat := range a.Entities.Thermostats {
		handler := func(state models.SensorState) {
			a.TimeSeriesDataStore.WriteTemperature(thermostat.ZoneID, thermostat.ID, state.DS18B20.Temperature)
		}

		a.MessageClient.SubscribeToThermostat(thermostat.Topic, handler)
	}
}

func (a *App) subscribeValves() {
	topicToHandlers := make(map[string][]func(map[string]any))
	for _, valve := range a.Entities.Valves {
		handler := func(state map[string]interface{}) {
			a.TimeSeriesDataStore.WriteValveState(valve.ZoneID, valve.ID, utils.ConvertOnOffToInt(state[valve.RelayName].(string)))
		}
		topicToHandlers[valve.StateTopic] = append(topicToHandlers[valve.StateTopic], handler)
	}

	for key, value := range topicToHandlers {
		log.Printf("Registering Topic: %s, Handlers: %d", key, len(value))
		a.MessageClient.SubscribeToValves(key, func(state map[string]interface{}) {
			for _, handler := range value {
				handler(state)
			}
		})
	}
}
