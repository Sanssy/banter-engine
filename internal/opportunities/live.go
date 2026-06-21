package opportunities

import (
	"fmt"

	"github.com/Sanssy/banter-engine/internal/matches"
)

func DetectLiveUpdates(previous, current []matches.Match) []Opportunity {
	previousByID := make(map[string]matches.Match, len(previous))
	for _, match := range previous {
		previousByID[match.MatchID] = match
	}

	var detected []Opportunity
	for _, match := range current {
		old, existed := previousByID[match.MatchID]
		if !existed {
			continue
		}

		label := fmt.Sprintf("%s - %s", match.HomeTeam, match.AwayTeam)
		if !isLiveStatus(old.Status) && isLiveStatus(match.Status) {
			detected = append(detected, Opportunity{Type: MatchStarted, Actor: label})
		}
		if old.Status != "fullTime" && match.Status == "fullTime" {
			detected = append(detected, Opportunity{Type: MatchEnded, Actor: label})
		}
		if old.Score != match.Score {
			detected = append(detected, Opportunity{
				Type:   ScoreChanged,
				Actor:  label,
				Target: fmt.Sprintf("%d-%d", match.Score.Home, match.Score.Away),
			})
			detected = append(detected, detectScoreEvents(old, match)...)
		}

		knownEvents := make(map[string]struct{}, len(old.Events))
		for _, event := range old.Events {
			knownEvents[eventKey(event)] = struct{}{}
		}
		for _, event := range match.Events {
			if !isImportantEvent(event.Type) {
				continue
			}
			if _, known := knownEvents[eventKey(event)]; known {
				continue
			}
			detected = append(detected, Opportunity{
				Type:   ImportantMatchEvent,
				Actor:  label,
				Target: fmt.Sprintf("%s %s", event.Type, event.Time),
			})
		}
	}

	return detected
}

func isLiveStatus(status string) bool {
	return status != "" && status != "preMatch" && status != "fullTime"
}

func isImportantEvent(eventType string) bool {
	switch eventType {
	case "goal", "booking", "penalty", "redCard", "var":
		return true
	default:
		return false
	}
}

func eventKey(event matches.Event) string {
	if event.ID != "" {
		return event.ID
	}
	return fmt.Sprintf(
		"%s|%s|%s|%s|%s|%d-%d",
		event.Type,
		event.Detail,
		event.Time,
		event.Side,
		event.PlayerID,
		event.Score.Home,
		event.Score.Away,
	)
}
