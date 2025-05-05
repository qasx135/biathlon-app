package models

import "time"

// Competitor an abstract of one competitor
type Competitor struct {
	ID              int
	Registered      bool
	StartSched      time.Time
	StartActual     time.Time
	Finished        bool
	CancelledFinish bool
	Disqualified    bool
	Comment         string
	LapTimes        []time.Duration
	LapStart        time.Time
	PenaltyTime     time.Duration
	PenaltyStart    time.Time
	Hits            int
	Shots           int
	CurrLap         int
	OnFiring        bool
	OnPenalty       bool
	LastEventTime   time.Time
}
