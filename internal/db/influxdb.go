package db

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	influxAPI "github.com/influxdata/influxdb-client-go/v2/api"
	"golang.org/x/net/context"
)

type InfluxDBClient struct {
	dbConnOnce sync.Once
	client     influxdb2.Client
	WriteAPI   influxAPI.WriteAPI
	QueryAPI   influxAPI.QueryAPI
}

type TimeSeriesDataStore interface {
	Close()
	WriteTemperature(zoneID int, thermostatID string, temp float64)
	WriteTemperatureTarget(zoneID int, target, error, output float64)
	WriteValveState(zoneID int, valveID int, state int)
	FetchZoneTemp(zoneID int, ctx context.Context) (float64, error)
}

func CreateInfluxClient() *InfluxDBClient {
	c := InfluxDBClient{}
	c.dbConnOnce.Do(func() {
		c.client = influxdb2.NewClientWithOptions("http://localhost:8086", "my-token", influxdb2.DefaultOptions().SetBatchSize(20))
		c.WriteAPI = c.client.WriteAPI("myorg", "mybucket")
		c.QueryAPI = c.client.QueryAPI("myorg")
	})
	return &c
}

func (c *InfluxDBClient) Close() {
	if c.WriteAPI != nil {
		c.WriteAPI.Flush()
	}
	if c.client != nil {
		c.client.Close()
	}
	log.Println("InfluxDB client closed")
}

func (c *InfluxDBClient) WriteTemperature(zoneID int, thermostatID string, temp float64) {
	p := influxdb2.NewPoint("temperature",
		map[string]string{"zoneID": strconv.Itoa(zoneID), "thermostatID": thermostatID},
		map[string]interface{}{"value": temp},
		time.Now())
	c.WriteAPI.WritePoint(p)
}

func (c *InfluxDBClient) WriteTemperatureTarget(zoneID int, target, error, output float64) {
	p := influxdb2.NewPoint("temperature",
		map[string]string{"zoneID": strconv.Itoa(zoneID)},
		map[string]interface{}{"error": error, "target": target, "output": output},
		time.Now())
	c.WriteAPI.WritePoint(p)
}

func (c *InfluxDBClient) WriteValveState(zoneID int, valveID int, state int) {
	p := influxdb2.NewPoint("valve_state",
		map[string]string{"zoneID": strconv.Itoa(zoneID), "valveID": strconv.Itoa(valveID)},
		map[string]interface{}{"value": state},
		time.Now())
	c.WriteAPI.WritePoint(p)
}

func (c *InfluxDBClient) FetchZoneTemp(zoneID int, ctx context.Context) (float64, error) {
	fluxQuery := fmt.Sprintf(`
    from(bucket: "mybucket")
  |> range(start: -1m)
  |> filter(fn: (r) => r._measurement == "temperature" and r._field == "value" and r.zoneID == "%d")
  |> mean()
  |> yield(name: "mean")
    `, zoneID)

	result, err := c.QueryAPI.Query(context.Background(), fluxQuery)
	if err != nil {
		log.Printf("Error querying InfluxDB: %s", err)
		return 0, err
	}

	if result.Next() {
		return result.Record().Value().(float64), nil
	} else if result.Err() != nil {
		log.Printf("Error reading result: %s", result.Err())
		return 0, result.Err()
	} else {
		log.Printf("No data for zoneID: %d", zoneID)
		return 0, fmt.Errorf("no data for zoneID: %d", zoneID)
	}
}
