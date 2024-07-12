package models

type Thermostat struct {
	ID     string `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	FriendlyName string `json:"friendly_name" db:"friendly_name"`
	Topic  string `json:"topic" db:"topic"`
	ZoneID int `json:"zone" db:"zone_id"`
}
