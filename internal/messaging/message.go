package messaging

import "github.com/monishth/go-therm/internal/models"

type MessageClient interface {
	Close()
	SendMessage(topic string, message string)
	SubscribeToThermostat(topic string, handler func(models.SensorState))
	SubscribeToValves(topic string, handler func(map[string]any))
}
