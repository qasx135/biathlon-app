package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"sunny_5_skiers/internal/config"
	"sunny_5_skiers/internal/models"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	fileConfig := os.Getenv("FILE_CONFIG")
	fileEvents := os.Getenv("FILE_EVENTS")
	configMain, err := config.LoadConfig(fileConfig)
	if err != nil {
		log.Fatal(err)
	}
	events := models.LoadEvents(fileEvents)

	fullCompetition := models.NewFullCompetition(configMain, events)
	fullCompetition.StartTime, err = time.Parse("15:04:05.000", configMain.Start)
	fullCompetition.StartDelta = MakeDuration(configMain.StartDelta)
	fullCompetition.Start()
	fullCompetition.GenerateOutput()
}

func MakeDuration(str string) time.Duration {
	t, err := time.Parse("15:04:05", str)
	if err != nil {
		log.Fatal(err)
	}
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second
}
