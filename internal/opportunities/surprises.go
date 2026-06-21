package opportunities

import (
	"fmt"

	"github.com/DSanoussy/banter-engine/internal/forecasts"
	"github.com/DSanoussy/banter-engine/internal/matches"
)

const (
	HugeUpset          = "HugeUpset"
	EveryoneWasWrong   = "EveryoneWasWrong"
	TheChosenOne       = "TheChosenOne"
	PredictionMassacre = "PredictionMassacre"
)

const (
	homeOutcome = "home"
	drawOutcome = "draw"
	awayOutcome = "away"
)

func DetectSurprises(match matches.Match, forecasts []forecasts.Forecast) []Opportunity {
	if match.Status != "fullTime" {
		return nil
	}

	actualOutcome := scoreOutcome(match.Score)
	matchLabel := fmt.Sprintf("%s - %s", match.HomeTeam, match.AwayTeam)
	correctShare := outcomeValue(match.PredictionStats, actualOutcome)
	statsAvailable := hasPredictionStats(match.PredictionStats)

	var detected []Opportunity
	if favorite, ok := favoriteOutcome(match.Quotations); ok && favorite != actualOutcome {
		detected = append(detected, Opportunity{
			Type:   HugeUpset,
			Actor:  outcomeName(match, actualOutcome),
			Target: outcomeName(match, favorite),
		})
	}
	detected = append(detected, DetectMassFailures(match)...)
	if statsAvailable && correctShare > 0 && correctShare < 0.05 {
		for _, forecast := range forecasts {
			if scoreOutcome(forecast.Prediction) == actualOutcome {
				detected = append(detected, Opportunity{
					Type:   TheChosenOne,
					Actor:  forecast.UserID,
					Target: matchLabel,
				})
			}
		}
	}
	return detected
}

func scoreOutcome(score matches.Score) string {
	switch {
	case score.Home > score.Away:
		return homeOutcome
	case score.Away > score.Home:
		return awayOutcome
	default:
		return drawOutcome
	}
}

func mostPredicted(stats matches.PredictionStats) (string, float64) {
	outcome := homeOutcome
	share := stats.Home
	if stats.Draw > share {
		outcome = drawOutcome
		share = stats.Draw
	}
	if stats.Away > share {
		outcome = awayOutcome
		share = stats.Away
	}
	return outcome, share
}

func outcomeValue(stats matches.PredictionStats, outcome string) float64 {
	switch outcome {
	case homeOutcome:
		return stats.Home
	case awayOutcome:
		return stats.Away
	default:
		return stats.Draw
	}
}

func favoriteOutcome(quotations matches.Quotations) (string, bool) {
	if quotations.Home <= 0 || quotations.Away <= 0 || quotations.Home == quotations.Away {
		return "", false
	}
	if quotations.Home < quotations.Away {
		return homeOutcome, true
	}
	return awayOutcome, true
}

func outcomeName(match matches.Match, outcome string) string {
	switch outcome {
	case homeOutcome:
		return match.HomeTeam
	case awayOutcome:
		return match.AwayTeam
	default:
		return "Match nul"
	}
}
