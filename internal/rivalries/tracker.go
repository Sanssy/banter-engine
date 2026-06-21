package rivalries

import (
	"sort"

	"github.com/Sanssy/banter-engine/internal/model"
	"github.com/Sanssy/banter-engine/internal/opportunities"
)

const rivalryThreshold = 3

type Rivalry struct {
	PlayerAID    string
	PlayerAName  string
	PlayerBID    string
	PlayerBName  string
	WinsA        int
	WinsB        int
	LastWinnerID string
	LastARank    int
	LastBRank    int
	LastAPoints  int
	LastBPoints  int
}

func Update(current []model.Standing, state []Rivalry) ([]Rivalry, []opportunities.Opportunity) {
	updated := append([]Rivalry(nil), state...)
	byPlayers := make(map[string]int, len(updated))
	for i := range updated {
		byPlayers[pairKey(updated[i].PlayerAID, updated[i].PlayerBID)] = i
	}

	var detected []opportunities.Opportunity
	for i := 0; i < len(current); i++ {
		for j := i + 1; j < len(current); j++ {
			playerA, playerB := orderedPlayers(current[i], current[j])
			key := pairKey(playerA.UserID, playerB.UserID)
			index, exists := byPlayers[key]
			if !exists {
				rivalry := newRivalry(playerA, playerB)
				updated = append(updated, rivalry)
				byPlayers[key] = len(updated) - 1
				continue
			}

			rivalry := &updated[index]
			if unchanged(*rivalry, playerA, playerB) || playerA.Rank == playerB.Rank {
				continue
			}

			oldDifference := abs(rivalry.WinsA - rivalry.WinsB)
			winner, loser := playerA, playerB
			if playerB.Rank < playerA.Rank {
				winner, loser = playerB, playerA
			}

			if rivalry.LastWinnerID != "" && rivalry.LastWinnerID != winner.UserID {
				detected = append(detected, opportunities.Opportunity{
					Type:   opportunities.Revenge,
					Actor:  winner.Name,
					Target: loser.Name,
				})
			}

			oldWinnerScore := rivalry.WinsB
			if winner.UserID == rivalry.PlayerAID {
				oldWinnerScore = rivalry.WinsA
				rivalry.WinsA++
			} else {
				rivalry.WinsB++
			}
			if oldWinnerScore < rivalryThreshold && oldWinnerScore+1 >= rivalryThreshold {
				detected = append(detected, opportunities.Opportunity{
					Type:   opportunities.Nemesis,
					Actor:  winner.Name,
					Target: loser.Name,
				})
			}
			if oldDifference < rivalryThreshold && abs(rivalry.WinsA-rivalry.WinsB) >= rivalryThreshold {
				detected = append(detected, opportunities.Opportunity{
					Type:   opportunities.Dominance,
					Actor:  winner.Name,
					Target: loser.Name,
				})
			}

			rivalry.LastWinnerID = winner.UserID
			setLastComparison(rivalry, playerA, playerB)
		}
	}

	sort.Slice(updated, func(i, j int) bool {
		return pairKey(updated[i].PlayerAID, updated[i].PlayerBID) < pairKey(updated[j].PlayerAID, updated[j].PlayerBID)
	})
	return updated, detected
}

func newRivalry(playerA, playerB model.Standing) Rivalry {
	rivalry := Rivalry{
		PlayerAID:   playerA.UserID,
		PlayerAName: playerA.Name,
		PlayerBID:   playerB.UserID,
		PlayerBName: playerB.Name,
	}
	if playerA.Rank < playerB.Rank {
		rivalry.WinsA = 1
		rivalry.LastWinnerID = playerA.UserID
	} else if playerB.Rank < playerA.Rank {
		rivalry.WinsB = 1
		rivalry.LastWinnerID = playerB.UserID
	}
	setLastComparison(&rivalry, playerA, playerB)
	return rivalry
}

func orderedPlayers(first, second model.Standing) (model.Standing, model.Standing) {
	if first.UserID < second.UserID {
		return first, second
	}
	return second, first
}

func pairKey(firstID, secondID string) string {
	if firstID < secondID {
		return firstID + "|" + secondID
	}
	return secondID + "|" + firstID
}

func unchanged(rivalry Rivalry, playerA, playerB model.Standing) bool {
	return rivalry.LastARank == playerA.Rank &&
		rivalry.LastBRank == playerB.Rank &&
		rivalry.LastAPoints == playerA.Points &&
		rivalry.LastBPoints == playerB.Points
}

func setLastComparison(rivalry *Rivalry, playerA, playerB model.Standing) {
	rivalry.PlayerAName = playerA.Name
	rivalry.PlayerBName = playerB.Name
	rivalry.LastARank = playerA.Rank
	rivalry.LastBRank = playerB.Rank
	rivalry.LastAPoints = playerA.Points
	rivalry.LastBPoints = playerB.Points
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
