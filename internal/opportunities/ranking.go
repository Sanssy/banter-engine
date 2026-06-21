package opportunities

import "github.com/DSanoussy/banter-engine/internal/model"

const RankingOvertake = "RankingOvertake"

type Opportunity struct {
	Type   string
	Actor  string
	Target string
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
