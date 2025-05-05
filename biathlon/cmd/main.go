package main

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"sunny_5_skiers/internal/config"
	"sunny_5_skiers/internal/models"
	"time"
)

func main() {
	// Load files from env
	err := godotenv.Load(".env")
	fileConfig := os.Getenv("FILE_CONFIG")
	fileEvents := os.Getenv("FILE_EVENTS")

	// Load config
	configMain, err := config.LoadConfig(fileConfig)
	if err != nil {
		slog.Error("Error loading config: ", err)
	}
	events := models.LoadEvents(fileEvents)

	fullCompetition := models.NewFullCompetition(configMain, events)
	fullCompetition.StartTime, err = time.Parse("15:04:05.000", configMain.Start)
	if err != nil {
		slog.Error("Error parsing start time: ", err)
	}
	fullCompetition.StartDelta = MakeDuration(configMain.StartDelta)
	fullCompetition.Start()
	fullCompetition.GenerateOutput()
}

func MakeDuration(str string) time.Duration {
	t, err := time.Parse("15:04:05", str)
	if err != nil {
		slog.Error("Error parsing time: ", err)
	}
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second
}
