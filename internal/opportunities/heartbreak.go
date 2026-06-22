package opportunities

import (
	"fmt"
	"strings"

	"github.com/Sanssy/banter-engine/internal/forecasts"
	"github.com/Sanssy/banter-engine/internal/matches"
)

func DetectHeartbreaks(previous, current matches.Match, forecasts []forecasts.Forecast) []Opportunity {
	matchLabel := fmt.Sprintf("%s - %s", current.HomeTeam, current.AwayTeam)
	var detected []Opportunity
	for _, event := range newMatchEvents(previous, current) {
		for _, forecast := range forecasts {
			switch {
			case isLateGoal(event) && forecast.Prediction == previous.Score && forecast.Prediction != current.Score:
				detected = append(detected, Opportunity{
					Type:    NinetiethMinuteHeartbreak,
					Actor:   forecastUserName(forecast),
					Target:  matchLabel,
					MatchID: current.MatchID,
					EventID: eventKey(event),
				})
			case isVAR(event) && forecast.Prediction == previous.Score && forecast.Prediction != current.Score:
				detected = append(detected, Opportunity{
					Type:    VARVictim,
					Actor:   forecastUserName(forecast),
					Target:  matchLabel,
					MatchID: current.MatchID,
					EventID: eventKey(event),
				})
			case isRedCard(event) && predictedWinner(forecast.Prediction) == event.Side:
				detected = append(detected, Opportunity{
					Type:    RedCardDisaster,
					Actor:   forecastUserName(forecast),
					Target:  matchLabel,
					MatchID: current.MatchID,
					EventID: eventKey(event),
				})
			}
		}
	}
	return EnsureIdentities(detected)
}

func DetectLatePointImpacts(
	previousMatch, currentMatch matches.Match,
	previousForecasts, currentForecasts []forecasts.Forecast,
) []Opportunity {
	var lateEventID string
	var hasAddedTimeEvent bool
	for _, event := range newMatchEvents(previousMatch, currentMatch) {
		if !isLateEvent(event) {
			continue
		}
		lateEventID = eventKey(event)
		if strings.Contains(event.Time, "+") {
			hasAddedTimeEvent = true
		}
	}
	if lateEventID == "" {
		return nil
	}

	matchLabel := fmt.Sprintf("%s - %s", currentMatch.HomeTeam, currentMatch.AwayTeam)
	var detected []Opportunity
	for _, impact := range calculatePointImpacts(previousForecasts, currentForecasts) {
		if impact.MatchID != currentMatch.MatchID {
			continue
		}
		switch {
		case impact.Delta < 0 && hasAddedTimeEvent:
			detected = append(detected, Opportunity{
				Type:    AddedTimeDisaster,
				Actor:   impact.UserName,
				Target:  matchLabel,
				MatchID: currentMatch.MatchID,
				EventID: lateEventID,
			})
		case impact.Delta > 0:
			detected = append(detected, Opportunity{
				Type:    LastMinuteHero,
				Actor:   impact.UserName,
				Target:  matchLabel,
				MatchID: currentMatch.MatchID,
				EventID: lateEventID,
			})
		}
	}
	return EnsureIdentities(detected)
}

func newMatchEvents(previous, current matches.Match) []matches.Event {
	knownEvents := make(map[string]struct{}, len(previous.Events))
	for _, event := range previous.Events {
		knownEvents[eventKey(event)] = struct{}{}
	}

	var events []matches.Event
	for _, event := range current.Events {
		if _, known := knownEvents[eventKey(event)]; !known {
			events = append(events, event)
		}
	}
	return events
}

func isLateGoal(event matches.Event) bool {
	if event.Type != "goal" || event.Detail == "var" {
		return false
	}
	var minute int
	_, err := fmt.Sscanf(event.Time, "%d'", &minute)
	return err == nil && minute >= 90
}

func isLateEvent(event matches.Event) bool {
	if event.Type != "goal" && !isVAR(event) {
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
