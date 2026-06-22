package opportunities

import (
	"fmt"
	"math"

	"github.com/Sanssy/banter-engine/internal/matches"
)

const crowdFavoriteThreshold = 0.60

func detectCrowdIntelligence(match matches.Match) []Opportunity {
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

func detectCrowdTransitions(previous, current matches.Match) []Opportunity {
	currentFavorite, currentFavoriteOK := crowdFavorite(current)
	previousFavorite, previousFavoriteOK := crowdFavorite(previous)

	var detected []Opportunity
	if currentFavoriteOK && (!previousFavoriteOK || currentFavorite.Actor != previousFavorite.Actor) {
		detected = append(detected, currentFavorite)
	}

	currentTrap, currentTrapOK := crowdTrap(current)
	previousTrap, previousTrapOK := crowdTrap(previous)
	if currentTrapOK && (!previousTrapOK || currentTrap.Actor != previousTrap.Actor || currentTrap.Target != previousTrap.Target) {
		detected = append(detected, currentTrap)
	}
	return detected
}

func crowdFavorite(match matches.Match) (Opportunity, bool) {
	if !hasPredictionStats(match.PredictionStats) {
		return Opportunity{}, false
	}
	popularOutcome, popularShare := mostPredicted(match.PredictionStats)
	if popularShare < crowdFavoriteThreshold {
		return Opportunity{}, false
	}
	return Opportunity{
		Type:   CrowdFavorite,
		Actor:  outcomeName(match, popularOutcome),
		Target: fmt.Sprintf("%d%%", int(math.Round(popularShare*100))),
	}, true
}

func crowdTrap(match matches.Match) (Opportunity, bool) {
	favorite, ok := crowdFavorite(match)
	if !ok {
		return Opportunity{}, false
	}
	popularOutcome, _ := mostPredicted(match.PredictionStats)
	bettingFavorite, ok := favoriteOutcome(match.Quotations)
	if !ok || bettingFavorite == popularOutcome {
		return Opportunity{}, false
	}
	return Opportunity{
		Type:   CrowdTrap,
		Actor:  favorite.Actor,
		Target: outcomeName(match, bettingFavorite),
	}, true
}

func detectPopularMistake(match matches.Match) []Opportunity {
	if !hasPredictionStats(match.PredictionStats) {
		return nil
	}
	popularOutcome, popularShare := mostPredicted(match.PredictionStats)
	if popularShare < crowdFavoriteThreshold || scoreOutcome(match.Score) == popularOutcome {
		return nil
	}
	return []Opportunity{{
		Type:   PopularMistake,
		Actor:  outcomeName(match, popularOutcome),
		Target: fmt.Sprintf("%s - %s", match.HomeTeam, match.AwayTeam),
	}}
}
