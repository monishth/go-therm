package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	InfluxDBURL    string `mapstructure:"INFLUXDB_URL"`
	InfluxDBToken  string `mapstructure:"INFLUXDB_TOKEN"`
	InfluxDBOrg    string `mapstructure:"INFLUXDB_ORG"`
	InfluxDBBucket string `mapstructure:"INFLUXDB_BUCKET"`
	InfluxDBPort   string `mapstructure:"INFLUXDB_PORT"`

	DBFile string `mapstructure:"DB_FILE"`

	Port string `mapstructure:"PORT"`

	MqttURL      string `mapstructure:"MQTT_URL"`
	MqttPort     string `mapstructure:"MQTT_PORT"`
	MqttClientID string `mapstructure:"MQTT_CLIENT_ID"`
}

func LoadConfig() *Config {
	// Defaults
	viper.SetDefault("INFLUXDB_URL", "localhost")
	viper.SetDefault("INFLUXDB_TOKEN", "my-token")
	viper.SetDefault("INFLUXDB_ORG", "myorg")
	viper.SetDefault("INFLUXDB_BUCKET", "mybucket")
	viper.SetDefault("INFLUXDB_PORT", "8086")
	viper.SetDefault("DB_FILE", "./config.db")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("MQTT_URL", "localhost")
	viper.SetDefault("MQTT_PORT", "1883")
	viper.SetDefault("MQTT_CLIENT_ID", "go-therm")

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	config := Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("error loading config")
	}

	return &config
}
