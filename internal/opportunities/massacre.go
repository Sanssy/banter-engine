package opportunities

import (
	"fmt"

	"github.com/DSanoussy/banter-engine/internal/matches"
)

func detectMassFailures(match matches.Match) []Opportunity {
	if match.Status != "fullTime" || !hasPredictionStats(match.PredictionStats) {
		return nil
	}

	actualOutcome := scoreOutcome(match.Score)
	correctShare := outcomeValue(match.PredictionStats, actualOutcome)
	mostPredictedOutcome, mostPredictedShare := mostPredicted(match.PredictionStats)
	matchLabel := fmt.Sprintf("%s - %s", match.HomeTeam, match.AwayTeam)

	var detected []Opportunity
	if mostPredictedShare > 0.80 && mostPredictedOutcome != actualOutcome {
		detected = append(detected, Opportunity{Type: EveryoneWasWrong, Actor: matchLabel})
	}
	if 1-correctShare > 0.50 {
		detected = append(detected, Opportunity{Type: PredictionMassacre, Actor: matchLabel})
	}
	return detected
}

func hasPredictionStats(stats matches.PredictionStats) bool {
	return stats.Home+stats.Draw+stats.Away > 0
}
