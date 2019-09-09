package config

import (
	"encoding/json"
	"log"
	"os"
)

var (
	host  string
	dbURI string
)

type Config struct {
	HostName     string
	Port         string
	DatabaseHost string
	RedisHost    string
}

func Init(environment string) *Config {
	file, _ := os.Open(`config/` + environment + `.json`)
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}

	err := decoder.Decode(&config)

	if err != nil {
		log.Fatal("Failed to read the ")
	}
	return config
}
