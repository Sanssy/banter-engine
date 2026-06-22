package opportunities

import (
	"fmt"

	"github.com/Sanssy/banter-engine/internal/matches"
)

func detectScoreEvents(previous, current matches.Match) []Opportunity {
	if previous.Score == current.Score {
		return nil
	}

	matchLabel := fmt.Sprintf("%s - %s", current.HomeTeam, current.AwayTeam)
	scoringTeam := scoreChangeActor(previous, current)
	detected := []Opportunity{{
		Type:    GoalSwing,
		Actor:   scoringTeam,
		Target:  fmt.Sprintf("%s (%d-%d)", matchLabel, current.Score.Home, current.Score.Away),
		MatchID: current.MatchID,
	}}

	previousOutcome := scoreOutcome(previous.Score)
	currentOutcome := scoreOutcome(current.Score)
	if previousOutcome == currentOutcome {
		return detected
	}
	if currentOutcome == drawOutcome {
		return append(detected, Opportunity{
			Type:    EqualizerChaos,
			Actor:   scoringTeam,
			Target:  matchLabel,
			MatchID: current.MatchID,
		})
	}
	if previousOutcome != drawOutcome {
		detected = append(detected, Opportunity{
			Type:    MatchTurnaround,
			Actor:   outcomeName(current, currentOutcome),
			Target:  outcomeName(previous, previousOutcome),
			MatchID: current.MatchID,
		})
	}
	return detected
}

func scoreChangeActor(previous, current matches.Match) string {
	switch {
	case current.Score.Home > previous.Score.Home && current.Score.Away == previous.Score.Away:
		return current.HomeTeam
	case current.Score.Away > previous.Score.Away && current.Score.Home == previous.Score.Home:
		return current.AwayTeam
	default:
		return fmt.Sprintf("%s - %s", current.HomeTeam, current.AwayTeam)
	}
}
