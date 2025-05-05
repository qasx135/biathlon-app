package models

import (
	"bufio"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// Event -> one event from line from file with events
type Event struct {
	Time             time.Time
	EventDI          int
	CompetitorsID    int
	AdditionalParams string
}

// LoadEvents takes string fileName gives an array of Event pointers
func LoadEvents(fileName string) []*Event {
	file, err := os.Open(fileName)
	if err != nil {
		slog.Error("Error opening file: ", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			slog.Error("Error closing file: ", err)
		}
	}(file)
	var events []*Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		event := parseLineToEvent(line)
		events = append(events, event)
	}
	return events
}

// parseLineToEvent takes string line gives a pointer on Event
func parseLineToEvent(line string) *Event {
	// split one line to separated parts
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
	//check if there is additional information
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
