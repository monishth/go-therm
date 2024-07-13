package models

type Valve struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	FriendlyName string `json:"friendly_name" db:"friendly_name"`
	StateTopic   string `json:"state_topic" db:"state_topic"`
	CommandTopic string `json:"command_topic" db:"command_topic"`
	RelayName    string `json:"relay_name" db:"relay_name"`
	ZoneID       int    `json:"zone_id" db:"zone_id"`
}
