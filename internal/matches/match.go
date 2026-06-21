package matches

type Match struct {
	MatchID         string
	HomeTeam        string
	AwayTeam        string
	Score           Score
	Status          string
	Quotations      Quotations
	PredictionStats PredictionStats
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
