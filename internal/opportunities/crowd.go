package opportunities

import (
	"fmt"
	"math"

	"github.com/DSanoussy/banter-engine/internal/matches"
)

const crowdFavoriteThreshold = 0.60

func DetectCrowdIntelligence(match matches.Match) []Opportunity {
	if !hasPredictionStats(match.PredictionStats) {
		return nil
	}

	popularOutcome, popularShare := mostPredicted(match.PredictionStats)
	if popularShare < crowdFavoriteThreshold {
		return nil
	}

	popularName := outcomeName(match, popularOutcome)
	detected := []Opportunity{{
		Type:   CrowdFavorite,
		Actor:  popularName,
		Target: fmt.Sprintf("%d%%", int(math.Round(popularShare*100))),
	}}

	if bettingFavorite, ok := favoriteOutcome(match.Quotations); ok && bettingFavorite != popularOutcome {
		detected = append(detected, Opportunity{
			Type:   CrowdTrap,
			Actor:  popularName,
			Target: outcomeName(match, bettingFavorite),
		})
	}
	if match.Status == "fullTime" && scoreOutcome(match.Score) != popularOutcome {
		detected = append(detected, Opportunity{
			Type:   PopularMistake,
			Actor:  popularName,
			Target: fmt.Sprintf("%s - %s", match.HomeTeam, match.AwayTeam),
		})
	}
	return detected
}
