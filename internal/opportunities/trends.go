package opportunities

import (
	"strconv"

	"github.com/Sanssy/banter-engine/internal/model"
)

const (
	minimumRankSwing         = 3
	minimumLeaderGapIncrease = 10
)

func detectRankingTrends(previous, current []model.Standing) []Opportunity {
	previousByUserID := make(map[string]model.Standing, len(previous))
	for _, standing := range previous {
		previousByUserID[standing.UserID] = standing
	}

	var detected []Opportunity
	for _, standing := range current {
		old, existed := previousByUserID[standing.UserID]
		if !existed {
			continue
		}

		change := old.Rank - standing.Rank
		switch {
		case change >= minimumRankSwing:
			detected = append(detected, Opportunity{
				Type:   ComebackSeason,
				Actor:  standing.Name,
				Target: strconv.Itoa(change),
			})
		case change <= -minimumRankSwing:
			detected = append(detected, Opportunity{
				Type:   FreeFall,
				Actor:  standing.Name,
				Target: strconv.Itoa(-change),
			})
		}
	}

	previousLeader, hasPreviousLeader := standingAtRank(previous, 1)
	previousSecond, hasPreviousSecond := standingAtRank(previous, 2)
	currentLeader, hasCurrentLeader := standingAtRank(current, 1)
	currentSecond, hasCurrentSecond := standingAtRank(current, 2)
	if hasPreviousLeader && hasPreviousSecond && hasCurrentLeader && hasCurrentSecond &&
		previousLeader.UserID == currentLeader.UserID {
		previousGap := previousLeader.Points - previousSecond.Points
		currentGap := currentLeader.Points - currentSecond.Points
		if currentGap-previousGap >= minimumLeaderGapIncrease {
			detected = append(detected, Opportunity{
				Type:   RunawayLeader,
				Actor:  currentLeader.Name,
				Target: currentSecond.Name,
			})
		}
	}

	return detected
}
