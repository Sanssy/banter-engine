package opportunities

import "github.com/DSanoussy/banter-engine/internal/model"

const maximumPodiumFightGap = 5

func detectPodiumFight(previous, current []model.Standing) []Opportunity {
	previousThird, hasPreviousThird := standingAtRank(previous, 3)
	previousFourth, hasPreviousFourth := standingAtRank(previous, 4)
	currentThird, hasCurrentThird := standingAtRank(current, 3)
	currentFourth, hasCurrentFourth := standingAtRank(current, 4)
	if !hasPreviousThird || !hasPreviousFourth || !hasCurrentThird || !hasCurrentFourth {
		return nil
	}

	previousGap := previousThird.Points - previousFourth.Points
	currentGap := currentThird.Points - currentFourth.Points
	playersChanged := previousThird.UserID != currentThird.UserID || previousFourth.UserID != currentFourth.UserID
	if currentGap > maximumPodiumFightGap || currentGap < 0 {
		return nil
	}
	if previousGap <= maximumPodiumFightGap && !playersChanged {
		return nil
	}

	return []Opportunity{{
		Type:   PodiumFight,
		Actor:  currentFourth.Name,
		Target: currentThird.Name,
	}}
}
