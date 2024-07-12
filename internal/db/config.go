package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/monishth/go-therm/internal/models"
)

// LoadConfig loads the configuration from the database
// This would usually be just the connection in this file
// but to keep things simple for now
// I do not want business logic in here permanently
type Config struct {
	Zones            []models.Zone
	Valves           []models.Valve
	Thermostats      []models.Thermostat
	ZoneToThermostat map[int][]*models.Thermostat
	ZoneToValve      map[int][]*models.Valve
	IdToZone         map[int]*models.Zone
}

func LoadConfig() Config {
	db, err := sqlx.Connect("sqlite3", "./config.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	zones := []models.Zone{}
	thermostats := []models.Thermostat{}
	valves := []models.Valve{}

	err = db.Select(&zones, "SELECT * FROM zone")
	if err != nil {
		panic(err)
	}

	err = db.Select(&thermostats, "SELECT * FROM thermostat")
	if err != nil {
		panic(err)
	}

	err = db.Select(&valves, "SELECT * FROM valve")
	if err != nil {
		panic(err)
	}
	zoneToThermostat := make(map[int][]*models.Thermostat)
	zoneToValve := make(map[int][]*models.Valve)

	for i := range thermostats {
		zoneToThermostat[thermostats[i].ZoneID] = append(zoneToThermostat[thermostats[i].ZoneID], &thermostats[i])
	}

	for i := range valves {
		zoneToValve[valves[i].ZoneID] = append(zoneToValve[valves[i].ZoneID], &valves[i])
	}

	idToZone := make(map[int]*models.Zone)
	for i := range zones {
		idToZone[zones[i].ID] = &zones[i]
	}

	return Config{
		Zones:            zones,
		Thermostats:      thermostats,
		Valves:           valves,
		ZoneToThermostat: zoneToThermostat,
		ZoneToValve:      zoneToValve,
		IdToZone:         idToZone,
	}
}
