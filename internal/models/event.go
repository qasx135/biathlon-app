package models

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Time             time.Time
	EventDI          int
	CompetitorsID    int
	AdditionalParams string
}

func LoadEvents(fileName string) []*Event {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var events []*Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		event := parseLineToEvent(line)
		events = append(events, event)
	}
	return events
}

func parseLineToEvent(line string) *Event {
	parts := strings.Split(line, " ")
	eventTimeString := strings.Trim(parts[0], "][")
	eventTimeUTC, err := time.Parse("15:04:05.000", eventTimeString)
	if err != nil {
		panic(err)
	}
	eventID, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	competitorsDI, err := strconv.Atoi(parts[2])
	if err != nil {
		panic(err)
	}
	var additionalParams string
	if len(parts) > 3 {
		additionalParams = parts[3]
	} else {
		additionalParams = ""
	}
	return &Event{Time: eventTimeUTC,
		EventDI:          eventID,
		CompetitorsID:    competitorsDI,
		AdditionalParams: additionalParams,
	}
}
