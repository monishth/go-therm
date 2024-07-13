package db

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	influxAPI "github.com/influxdata/influxdb-client-go/v2/api"
)

var (
	dbConnOnce sync.Once
	client     influxdb2.Client
	writeAPI   influxAPI.WriteAPI
)

func GetInfluxClient() influxAPI.WriteAPI {
	dbConnOnce.Do(func() {
		client = influxdb2.NewClientWithOptions("http://localhost:8086", "my-token", influxdb2.DefaultOptions().SetBatchSize(20))
		writeAPI = client.WriteAPI("myorg", "mybucket")
	})
	return writeAPI
}

func CloseInfluxClient() {
	if writeAPI != nil {
		writeAPI.Flush()
	}
	if client != nil {
		client.Close()
	}
	log.Println("InfluxDB client closed")
}

func WriteTemperature(writeAPI influxAPI.WriteAPI, zoneID int, thermostatID string, temp float64) {
	p := influxdb2.NewPoint("temperature",
		map[string]string{"zoneID": strconv.Itoa(zoneID), "thermostatID": thermostatID},
		map[string]interface{}{"value": temp},
		time.Now())
	writeAPI.WritePoint(p)
}

func WriteValveState(writeAPI influxAPI.WriteAPI, zoneID int, valveID int, state int) {
	p := influxdb2.NewPoint("valve_state",
		map[string]string{"zoneID": strconv.Itoa(zoneID), "valveID": strconv.Itoa(valveID)},
		map[string]interface{}{"value": state},
		time.Now())
	writeAPI.WritePoint(p)
}
