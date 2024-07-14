package db

import "context"

type TimeSeriesDataStore interface {
	Close()
	WriteTemperature(zoneID int, thermostatID string, temp float64)
	WriteTemperatureData(zoneID int, target, error, output float64)
	WriteTarget(zoneID int, target float64)
	WriteValveState(zoneID int, valveID int, state int)
	FetchZoneTemp(zoneID int, ctx context.Context) (float64, error)
}
