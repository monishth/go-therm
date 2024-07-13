package messaging

import (
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/monishth/go-therm/internal/models"
	"github.com/monishth/go-therm/pkg/utils"
)

type MQTTClient struct {
	client MQTT.Client
}

var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	payload := msg.Payload()
	log.Printf("Received message: %s from topic: %s\n", payload, msg.Topic())
}

func StartMQTTClient() MQTTClient {
	opts := MQTT.NewClientOptions().AddBroker("tcp://172.16.255.82:1883").SetClientID("go-therm")
	opts.SetDefaultPublishHandler(messagePubHandler)
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	log.Println("Connected to MQTT broker")

	return MQTTClient{client}
}

func (c *MQTTClient) Close() {
	c.client.Disconnect(250)
	log.Println("Disconnected from MQTT broker")
}

func (c *MQTTClient) SendMessage(topic string, message string) {
	token := c.client.Publish(topic, 0, false, message)
	token.Wait()
}

func (c *MQTTClient) subscribeToTopic(topic string, handler MQTT.MessageHandler) error {
	if token := c.client.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (c *MQTTClient) SubscribeToThermostat(topic string, handler func(models.SensorState)) {
	c.subscribeToTopic(topic, WrapHandler(handler))
}

func (c *MQTTClient) SubscribeToValves(topic string, handler func(map[string]any)) {
	c.subscribeToTopic(topic, WrapHandler(handler))
}

func WrapHandler[T any](handler func(T)) MQTT.MessageHandler {
	return func(client MQTT.Client, msg MQTT.Message) {
		payload := msg.Payload()
		value, err := utils.Decode[T](payload)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Received message: %s from topic: %s\n", payload, msg.Topic())
		handler(value)
	}
}
