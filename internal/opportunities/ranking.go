package opportunities

import "github.com/DSanoussy/banter-engine/internal/model"

const (
	RankingOvertake     = "RankingOvertake"
	EnteredTop3         = "EnteredTop3"
	ExitedTop3          = "ExitedTop3"
	LeaderUnderPressure = "LeaderUnderPressure"
	LastPlaceLocked     = "LastPlaceLocked"
	Nemesis             = "Nemesis"
	Revenge             = "Revenge"
	Dominance           = "Dominance"
)

type Opportunity struct {
	Type   string
	Actor  string
	Target string
}

func Detect(previous, current []model.Standing) []Opportunity {
	if len(previous) == 0 || len(current) == 0 {
		return nil
	}

	detected := DetectRankingOvertakes(previous, current)
	detected = append(detected, detectTop3Changes(previous, current)...)
	detected = append(detected, DetectRankingTrends(previous, current)...)

	if opportunity, ok := detectLeaderUnderPressure(previous, current); ok {
		detected = append(detected, opportunity)
	}
	if opportunity, ok := detectLastPlaceLocked(previous, current); ok {
		detected = append(detected, opportunity)
	}

	return detected
}

func DetectRankingOvertakes(previous, current []model.Standing) []Opportunity {
	previousByUserID := make(map[string]model.Standing, len(previous))
	for _, standing := range previous {
		previousByUserID[standing.UserID] = standing
	}

	var detected []Opportunity
	for _, actor := range current {
		previousActor, actorExisted := previousByUserID[actor.UserID]
		if !actorExisted {
			continue
		}

		for _, target := range current {
			if actor.UserID == target.UserID {
				continue
			}

			previousTarget, targetExisted := previousByUserID[target.UserID]
			if !targetExisted {
				continue
			}

			if previousActor.Rank > previousTarget.Rank && actor.Rank < target.Rank {
				detected = append(detected, Opportunity{
					Type:   RankingOvertake,
					Actor:  actor.Name,
					Target: target.Name,
				})
			}
		}
	}

	return detected
}

func detectTop3Changes(previous, current []model.Standing) []Opportunity {
	previousByUserID := make(map[string]model.Standing, len(previous))
	for _, standing := range previous {
		previousByUserID[standing.UserID] = standing
	}

	var detected []Opportunity
	for _, standing := range current {
		previousStanding, existed := previousByUserID[standing.UserID]
		if !existed {
			continue
		}

		switch {
		case previousStanding.Rank > 3 && standing.Rank > 0 && standing.Rank <= 3:
			detected = append(detected, Opportunity{Type: EnteredTop3, Actor: standing.Name})
		case previousStanding.Rank > 0 && previousStanding.Rank <= 3 && standing.Rank > 3:
			detected = append(detected, Opportunity{Type: ExitedTop3, Actor: standing.Name})
		}
	}

	return detected
}

func detectLeaderUnderPressure(previous, current []model.Standing) (Opportunity, bool) {
	previousLeader, hasPreviousLeader := standingAtRank(previous, 1)
	previousSecond, hasPreviousSecond := standingAtRank(previous, 2)
	currentLeader, hasCurrentLeader := standingAtRank(current, 1)
	currentSecond, hasCurrentSecond := standingAtRank(current, 2)
	if !hasPreviousLeader || !hasPreviousSecond || !hasCurrentLeader || !hasCurrentSecond {
		return Opportunity{}, false
	}
	if previousLeader.UserID != currentLeader.UserID {
		return Opportunity{}, false
	}

	previousGap := previousLeader.Points - previousSecond.Points
	currentGap := currentLeader.Points - currentSecond.Points
	if currentGap >= previousGap {
		return Opportunity{}, false
	}

	return Opportunity{
		Type:   LeaderUnderPressure,
		Actor:  currentLeader.Name,
		Target: currentSecond.Name,
	}, true
}

func detectLastPlaceLocked(previous, current []model.Standing) (Opportunity, bool) {
	previousLast, hasPreviousLast := lastStanding(previous)
	currentLast, hasCurrentLast := lastStanding(current)
	if !hasPreviousLast || !hasCurrentLast || previousLast.UserID != currentLast.UserID {
		return Opportunity{}, false
	}

	return Opportunity{Type: LastPlaceLocked, Actor: currentLast.Name}, true
}

func standingAtRank(standings []model.Standing, rank int) (model.Standing, bool) {
	for _, standing := range standings {
		if standing.Rank == rank {
			return standing, true
		}
	}
	return model.Standing{}, false
}

func lastStanding(standings []model.Standing) (model.Standing, bool) {
	if len(standings) == 0 {
		return model.Standing{}, false
	}

	last := standings[0]
	for _, standing := range standings[1:] {
		if standing.Rank > last.Rank {
			last = standing
		}
	}
	return last, true
}
