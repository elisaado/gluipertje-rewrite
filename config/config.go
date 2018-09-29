package config

import (
	"encoding/json"
	"os"
	"fmt"
)

type Config struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	ExternalURL string `json:"external_url"`
}

var C Config

func InitConfig() {
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&C); err != nil {
		fmt.Println("Parsing config file", err.Error())
	}
}
