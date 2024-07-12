package models

type DS18B20 struct {
	ID          string  `json:"Id"`
	Temperature float64 `json:"Temperature"`
}

type SensorState struct {
	DS18B20  DS18B20 `json:"DS18B20"`
	TempUnit string  `json:"TempUnit"`
}

type ValveState struct {
	Heap      int    `json:"Heap"`
	LoadAvg   int    `json:"LoadAvg"`
	MqttCount int    `json:"MqttCount"`
	Power1    string `json:"POWER1"`
	Power2    string `json:"POWER2"`
	Power3    string `json:"POWER3"`
	Power4    string `json:"POWER4"`
}

type TempMeasurement struct {
	Zone string  `json:"zone"`
	Temp float64 `json:"temp"`
	
}
