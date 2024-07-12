package models

type Zone struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	FriendlyName string `json:"friendly_name" db:"friendly_name"`
}
