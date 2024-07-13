package models

type ZoneModel struct {
	ID                 int     `json:"id" db:"id"`
	Name               string  `json:"name" db:"name"`
	FriendlyName       string  `json:"friendly_name" db:"friendly_name"`
	OverrideTargetTemp float64 `json:"override_target_temp" db:"override_target_temp"`
}

type Zone struct {
	ID           int           `json:"id" `
	Name         string        `json:"name"`
	FriendlyName string        `json:"friendly_name"`
	Valves       []*Valve      `json:"valves"`
	Thermostats  []*Thermostat `json:"thermostats"`
}
