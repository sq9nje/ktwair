package shared

import (
	"encoding/json"
	"io/ioutil"
)

// Default configuration file location
const DefaultConfigFile = "/etc/ktwair/config.json"

// Global variable for configuration
var GlobalConfig Config

// Configuration file structure
type Config struct {
	LogLevel string `json:"log_level"`

	KTWAir struct {
		BaseURL       string `json:"base_url"`
		Timeout       int    `json:"timeout"`
		QueryInterval int    `json:"query_interval"`
		StationID     int    `json:"station_id"`
	} `json:"ktwair"`

	MQTT struct {
		Hostname string `json:"hostname"`
		Port     int    `json:"port"`
		UseTLS   bool   `json:"tls"`
		User     string `json:"user"`
		Password string `json:"password"`
		Topic    string `json:"topic"`
	} `json:"mqtt"`
}

func ReadConfig(filename string) (Config, error) {
	if filename == "" {
		filename = DefaultConfigFile
	}

	rawJson, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = json.Unmarshal(rawJson, &config)
	GlobalConfig = config
	return config, err
}
