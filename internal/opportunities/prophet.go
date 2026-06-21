package opportunities

import (
	"fmt"

	"github.com/DSanoussy/banter-engine/internal/forecasts"
	"github.com/DSanoussy/banter-engine/internal/matches"
)

func DetectProphets(match matches.Match, matchForecasts []forecasts.Forecast) []Opportunity {
	if match.Status != "fullTime" || !hasPredictionStats(match.PredictionStats) {
		return nil
	}

	actualOutcome := scoreOutcome(match.Score)
	correctShare := outcomeValue(match.PredictionStats, actualOutcome)
	if correctShare <= 0 || correctShare >= 0.05 {
		return nil
	}

	var correct []forecasts.Forecast
	for _, forecast := range matchForecasts {
		if scoreOutcome(forecast.Prediction) == actualOutcome {
			correct = append(correct, forecast)
		}
	}
	if len(correct) == 0 {
		return nil
	}

	matchLabel := fmt.Sprintf("%s - %s", match.HomeTeam, match.AwayTeam)
	var detected []Opportunity
	if len(correct) == 1 {
		detected = append(detected, Opportunity{
			Type:   TheChosenOne,
			Actor:  correct[0].UserID,
			Target: matchLabel,
		})
	}
	for _, forecast := range correct {
		detected = append(detected, Opportunity{
			Type:   AgainstTheCrowd,
			Actor:  forecast.UserID,
			Target: matchLabel,
		})
	}
	return detected
}
