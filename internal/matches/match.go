package matches

import "time"

type Match struct {
	MatchID         string
	Date            time.Time
	HomeTeam        string
	AwayTeam        string
	Score           Score
	Status          string
	Quotations      Quotations
	PredictionStats PredictionStats
	Events          []Event
}

type Score struct {
	Home int
	Away int
}

type Quotations struct {
	Home int
	Draw int
	Away int
}

type PredictionStats struct {
	Home float64
	Draw float64
	Away float64
}

type Event struct {
	ID       string
	Type     string
	Detail   string
	Time     string
	Side     string
	PlayerID string
	Score    Score
}
