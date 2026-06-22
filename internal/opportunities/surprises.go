package opportunities

import (
	"github.com/Sanssy/banter-engine/internal/forecasts"
	"github.com/Sanssy/banter-engine/internal/matches"
)

const (
	homeOutcome = "home"
	drawOutcome = "draw"
	awayOutcome = "away"
)

func DetectSurprises(previous, current matches.Match, forecasts []forecasts.Forecast) []Opportunity {
	var detected []Opportunity
	if previous.Status != "fullTime" && current.Status == "fullTime" {
		actualOutcome := scoreOutcome(current.Score)
		if favorite, ok := favoriteOutcome(current.Quotations); ok && favorite != actualOutcome {
			detected = append(detected, Opportunity{
				Type:   HugeUpset,
				Actor:  outcomeName(current, actualOutcome),
				Target: outcomeName(current, favorite),
			})
		}
		detected = append(detected, detectMassFailures(current)...)
		detected = append(detected, detectCrowdTransitions(previous, current)...)
		detected = append(detected, detectPopularMistake(current)...)
		detected = append(detected, detectProphets(current, forecasts)...)
		return detected
	}
	return detectCrowdTransitions(previous, current)
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
