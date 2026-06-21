package forecasts

import "github.com/DSanoussy/banter-engine/internal/matches"

type Forecast struct {
	UserID     string
	MatchID    string
	Prediction matches.Score
	Result     matches.Score
	Points     int
}
