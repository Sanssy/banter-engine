package forecasts

import (
	"time"

	"github.com/DSanoussy/banter-engine/internal/matches"
)

type Forecast struct {
	UserID     string
	MatchID    string
	MatchDate  time.Time
	Prediction matches.Score
	Result     matches.Score
	Points     int
}
