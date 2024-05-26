package main

import (
	"encoding/json"
	"os"
)

var (
	configPath                 = "config.json"
	config     AutomoverConfig = AutomoverConfig{
		Watchlist: []AutomoverWatchedFolder{},
	}
)

type AutomoverConfig struct {
	Watchlist []AutomoverWatchedFolder `json:"watchlist"`
}

type AutomoverWatchedFolder struct {
	WatchPath       string `json:"watch_path"`
	WatchPattern    string `json:"watch_pattern"`
	DestinationPath string `json:"destination_path"`
}

func loadConfig() error {
	file, err := os.Open(configPath)
	if err != nil {
		return saveConfig()
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}

	return nil
}

func saveConfig() error {
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}
