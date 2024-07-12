package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/monishth/go-therm/pkg/utils"
	"log"
)

var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	payload := msg.Payload()
	log.Printf("Received message: %s from topic: %s\n", payload, msg.Topic())
}

func StartMQTTClient() MQTT.Client {
	opts := MQTT.NewClientOptions().AddBroker("tcp://172.16.255.82:1883").SetClientID("go-therm")
	opts.SetDefaultPublishHandler(messagePubHandler)
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	log.Println("Connected to MQTT broker")

	return client
}

func ShutdownMQTTClient(client MQTT.Client) {
	client.Disconnect(250)
	log.Println("Disconnected from MQTT broker")
}

func SubscribeToTopic(client MQTT.Client, topic string, handler MQTT.MessageHandler) error {
	if token := client.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
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
