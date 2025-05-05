package models

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"sunny_5_skiers/internal/config"
	"time"
)

// FullCompetition abstraction of full process
type FullCompetition struct {
	Config         *config.Config
	AllCompetitors map[int]*Competitor
	Events         []*Event
	StartTime      time.Time
	StartDelta     time.Duration
}

// NewFullCompetition takes config, array of pointers on Event and gives pointer on FullCompetition
func NewFullCompetition(config *config.Config, events []*Event) *FullCompetition {
	return &FullCompetition{
		Config:         config,
		AllCompetitors: make(map[int]*Competitor),
		Events:         events,
	}
}

// LoggingEvent takes t time.Time, msg string
// simple logger to console
func (c *FullCompetition) LoggingEvent(t time.Time, msg string) {
	fmt.Printf("[%s] %s\n", t.Format("15:04:05.000"), msg)
}

// Start method for starting process of file with events and make information about competitors
func (c *FullCompetition) Start() {
	for _, event := range c.Events {
		competitor := c.GetCompetitor(event.CompetitorsID)
		competitor.LastEventTime = event.Time

		//simple choice of process steps of competitors (maybe need to make concurrency)
		switch event.EventDI {
		case 1:
			competitor.Registered = true
			c.LoggingEvent(event.Time, fmt.Sprintf("The competitor(%d) registered", competitor.ID))

		case 2:
			startTime, err := StringToTime(event.AdditionalParams)
			if err != nil {
				log.Printf("Invalid start time for competitor %d: %v", competitor.ID, err)
				continue
			}
			competitor.StartSched = startTime
			c.LoggingEvent(event.Time, fmt.Sprintf("The start time for competitor(%d) was set by a draw to %s",
				competitor.ID, TimeToString(startTime)))

		case 3:
			c.LoggingEvent(event.Time, fmt.Sprintf("The competitor(%d) is on the start line", competitor.ID))

		case 4:
			if event.Time.After(competitor.StartSched.Add(c.StartDelta)) {
				competitor.Disqualified = true
				c.LoggingEvent(event.Time, fmt.Sprintf("The competitor(%d) is disqualified", competitor.ID))
			}
			competitor.StartActual = event.Time
			c.LoggingEvent(event.Time, fmt.Sprintf("The competitor(%d) has started", competitor.ID))

		case 5:
			competitor.OnFiring = true
			c.LoggingEvent(event.Time, fmt.Sprintf("The competitor(%d) is on the firing range(%s)",
				competitor.ID, event.AdditionalParams))

		case 6:
			shots, _ := strconv.Atoi(event.AdditionalParams)
			if competitor.Shots < shots {
				competitor.Shots = shots
			}
			competitor.Hits++
			c.LoggingEvent(event.Time, fmt.Sprintf("The target(%s) has been hit by competitor(%d)",
				event.AdditionalParams, competitor.ID))

		case 7:
			competitor.OnFiring = false
			c.LoggingEvent(event.Time, fmt.Sprintf("The competitor(%d) left the firing range", competitor.ID))

		case 8:
			competitor.OnPenalty = true
			competitor.PenaltyStart = event.Time
			c.LoggingEvent(event.Time, fmt.Sprintf("The competitor(%d) entered the penalty laps", competitor.ID))

		case 9:
			competitor.OnPenalty = false
			if !competitor.PenaltyStart.IsZero() {
				competitor.PenaltyTime += event.Time.Sub(competitor.PenaltyStart)
			}
			c.LoggingEvent(event.Time, fmt.Sprintf("The competitor(%d) left the penalty laps", competitor.ID))

		case 10:
			if competitor.LapStart.IsZero() {

				competitor.LapStart = competitor.StartActual
			}
			lapTime := event.Time.Sub(competitor.LapStart)
			competitor.LapTimes = append(competitor.LapTimes, lapTime)
			competitor.LapStart = event.Time
			competitor.CurrLap++
			c.LoggingEvent(event.Time, fmt.Sprintf("The competitor(%d) ended the main lap", competitor.ID))

		case 11:
			competitor.CancelledFinish = true
			competitor.Comment = event.AdditionalParams
			c.LoggingEvent(event.Time, fmt.Sprintf("The competitor(%d) can't continue: %s", competitor.ID, event.AdditionalParams))
		}
	}
}

// GetCompetitor takes id int and gives pointer on Competitor
func (c *FullCompetition) GetCompetitor(id int) *Competitor {
	if _, ok := c.AllCompetitors[id]; !ok {
		c.AllCompetitors[id] = &Competitor{
			ID:       id,
			LapTimes: make([]time.Duration, 0),
			CurrLap:  1,
		}
	}
	return c.AllCompetitors[id]
}

// StringToTime takes timeStr string and gives time.Time, error
// make string to time.Time
func StringToTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04:05.000", timeStr)
}

// TimeToString takes t time.Time and gives string
// make time.Time to string
func TimeToString(t time.Time) string {
	return t.Format("15:04:05.000")
}

// GenerateOutput method to make simple UI for output
// information puts in console
func (c *FullCompetition) GenerateOutput() {
	fmt.Println("\nResulting table:")

	var allCompetitors []*Competitor
	for _, competitor := range c.AllCompetitors {
		allCompetitors = append(allCompetitors, competitor)
	}

	// sorting with time and include not finished persons
	sort.Slice(allCompetitors, func(i, j int) bool {
		if allCompetitors[i].Disqualified {
			return false
		}
		if allCompetitors[j].Disqualified {
			return true
		}
		if allCompetitors[i].CancelledFinish && !allCompetitors[j].CancelledFinish {
			return false
		}
		if !allCompetitors[i].CancelledFinish && allCompetitors[j].CancelledFinish {
			return true
		}

		timeI := allCompetitors[i].LastEventTime.Sub(allCompetitors[i].StartActual)
		timeJ := allCompetitors[j].LastEventTime.Sub(allCompetitors[j].StartActual)
		return timeI < timeJ
	})

	for _, competitor := range allCompetitors {
		if competitor.Disqualified {
			fmt.Printf("[NotStarted] %d\n", competitor.ID)
			continue
		}

		if competitor.CancelledFinish {
			fmt.Printf("[NotFinished] %d ", competitor.ID)
		} else {
			totalTime := competitor.LastEventTime.Sub(competitor.StartActual)
			fmt.Printf("[%s] %d ", DurationToString(totalTime), competitor.ID)
		}
		var lapTimes []string
		for i, lapTime := range competitor.LapTimes {
			if i < c.Config.Laps {
				speed := float64(c.Config.LapLen) / lapTime.Seconds()
				lapTimes = append(lapTimes, fmt.Sprintf("{%s, %.3f}",
					DurationToString(lapTime), speed))
			}
		}
		fmt.Printf("%v ", lapTimes)

		if competitor.PenaltyTime > 0 {
			speed := float64(c.Config.PenaltyLen) / competitor.PenaltyTime.Seconds()
			fmt.Printf("{%s, %.3f} ",
				DurationToString(competitor.PenaltyTime), speed)
		} else {
			fmt.Printf("{,} ")
		}

		fmt.Printf("%d/%d\n", competitor.Hits, competitor.Shots*2)
	}
}

// DurationToString takes d time.Duration and gives string
// make time.Duration to string format
func DurationToString(d time.Duration) string {
	d = d.Round(time.Millisecond)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ms := d / time.Millisecond
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}
