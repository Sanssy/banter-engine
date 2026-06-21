package opportunities

import (
	"fmt"

	"github.com/DSanoussy/banter-engine/internal/forecasts"
	"github.com/DSanoussy/banter-engine/internal/matches"
)

func DetectHeartbreaks(previous, current matches.Match, forecasts []forecasts.Forecast) []Opportunity {
	knownEvents := make(map[string]struct{}, len(previous.Events))
	for _, event := range previous.Events {
		knownEvents[eventKey(event)] = struct{}{}
	}

	matchLabel := fmt.Sprintf("%s - %s", current.HomeTeam, current.AwayTeam)
	var detected []Opportunity
	for _, event := range current.Events {
		if _, known := knownEvents[eventKey(event)]; known {
			continue
		}

		for _, forecast := range forecasts {
			switch {
			case isLateGoal(event) && forecast.Prediction == previous.Score && forecast.Prediction != current.Score:
				detected = append(detected, Opportunity{
					Type:   NinetiethMinuteHeartbreak,
					Actor:  forecast.UserID,
					Target: matchLabel,
				})
			case isVAR(event) && forecast.Prediction == previous.Score && forecast.Prediction != current.Score:
				detected = append(detected, Opportunity{
					Type:   VARVictim,
					Actor:  forecast.UserID,
					Target: matchLabel,
				})
			case isRedCard(event) && predictedWinner(forecast.Prediction) == event.Side:
				detected = append(detected, Opportunity{
					Type:   RedCardDisaster,
					Actor:  forecast.UserID,
					Target: matchLabel,
				})
			}
		}
	}
	return detected
}

func isLateGoal(event matches.Event) bool {
	if event.Type != "goal" || event.Detail == "var" {
		return false
	}
	var minute int
	_, err := fmt.Sscanf(event.Time, "%d'", &minute)
	return err == nil && minute >= 90
}

func isVAR(event matches.Event) bool {
	return event.Type == "var" || event.Detail == "var"
}

func isRedCard(event matches.Event) bool {
	if event.Type == "redCard" {
		return true
	}
	if event.Type != "booking" {
		return false
	}
	switch event.Detail {
	case "red", "straightRed", "secondYellow":
		return true
	default:
		return false
	}
}

func predictedWinner(score matches.Score) string {
	switch {
	case score.Home > score.Away:
		return "home"
	case score.Away > score.Home:
		return "away"
	default:
		return "draw"
	}
}
