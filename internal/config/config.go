package config

import (
	"encoding/json"
	"log/slog"
	"os"
)

type Config struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

func LoadConfig(path string) (*Config, error) {
	var conf Config
	data, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Error reading config file: ", err)
		return nil, err
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		slog.Error("Error parsing config file: ", err)
	}
	return &conf, nil
}
