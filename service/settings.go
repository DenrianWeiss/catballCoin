package service

import (
	"encoding/json"
	"github.com/DenrianWeiss/catballCoin/model"
	"io/ioutil"
	"os"
)

const (
	configPath = "config.json"
)

var (
	GlobalConfig model.GlobalConfig
)

func settingsInit() {
	configFile, err := os.Open(configPath)
	if err != nil {
		panic("No config.")
	}
	configContent, err := ioutil.ReadAll(configFile)
	if err != nil {
		panic("Failed to read config")
	}

	configResult := model.GlobalConfig{}

	err = json.Unmarshal(configContent, &configResult)

	if err != nil {
		panic("Illegal config")
	}

	GlobalConfig = configResult
}
