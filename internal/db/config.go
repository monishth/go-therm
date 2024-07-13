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
	IDToZone         map[int]*models.ZoneModel
}

func LoadConfig() Config {
	db, err := sqlx.Connect("sqlite3", "./config.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	zoneModels := []models.ZoneModel{}
	thermostats := []models.Thermostat{}
	valves := []models.Valve{}

	err = db.Select(&zoneModels, "SELECT * FROM zone")
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

	idToZone := make(map[int]*models.ZoneModel)
	for i := range zoneModels {
		idToZone[zoneModels[i].ID] = &zoneModels[i]
	}
	zones := make([]models.Zone, len(zoneModels))
	for i, zoneModel := range zoneModels {
		zones[i] = models.Zone{
			ID:           zoneModel.ID,
			Name:         zoneModel.Name,
			FriendlyName: zoneModel.FriendlyName,
			Valves:       zoneToValve[zoneModel.ID],
			Thermostats:  zoneToThermostat[zoneModel.ID],
		}
	}

	return Config{
		Zones:       zones,
		Thermostats: thermostats,
		Valves:      valves,
		IDToZone:    idToZone,
	}
}
