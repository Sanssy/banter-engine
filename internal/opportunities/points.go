package opportunities

import (
	"fmt"
	"sort"

	"github.com/Sanssy/banter-engine/internal/forecasts"
)

const pointExplosionThreshold = 5

type pointImpact struct {
	UserID         string
	MatchID        string
	PreviousPoints int
	CurrentPoints  int
	Delta          int
}

func calculatePointImpacts(previous, current []forecasts.Forecast) []pointImpact {
	previousByForecast := make(map[string]forecasts.Forecast, len(previous))
	for _, forecast := range previous {
		previousByForecast[forecastKey(forecast)] = forecast
	}

	var impacts []pointImpact
	for _, forecast := range current {
		old, existed := previousByForecast[forecastKey(forecast)]
		if !existed || old.Points == forecast.Points {
			continue
		}
		impacts = append(impacts, pointImpact{
			UserID:         forecast.UserID,
			MatchID:        forecast.MatchID,
			PreviousPoints: old.Points,
			CurrentPoints:  forecast.Points,
			Delta:          forecast.Points - old.Points,
		})
	}
	sort.Slice(impacts, func(i, j int) bool {
		if impacts[i].UserID == impacts[j].UserID {
			return impacts[i].MatchID < impacts[j].MatchID
		}
		return impacts[i].UserID < impacts[j].UserID
	})
	return impacts
}

func DetectPointImpacts(previous, current []forecasts.Forecast) []Opportunity {
	impacts := calculatePointImpacts(previous, current)
	if len(impacts) == 0 {
		return nil
	}

	var biggestWinner *pointImpact
	var biggestLoser *pointImpact
	for i := range impacts {
		impact := &impacts[i]
		if impact.Delta > 0 && (biggestWinner == nil || impact.Delta > biggestWinner.Delta) {
			biggestWinner = impact
		}
		if impact.Delta < 0 && (biggestLoser == nil || impact.Delta < biggestLoser.Delta) {
			biggestLoser = impact
		}
	}

	var detected []Opportunity
	if biggestWinner != nil {
		detected = append(detected, Opportunity{
			Type:   BiggestWinner,
			Actor:  biggestWinner.UserID,
			Target: signedPoints(biggestWinner.Delta),
		})
	}
	if biggestLoser != nil {
		detected = append(detected, Opportunity{
			Type:   BiggestLoser,
			Actor:  biggestLoser.UserID,
			Target: signedPoints(biggestLoser.Delta),
		})
	}
	for _, impact := range impacts {
		if impact.Delta >= pointExplosionThreshold {
			detected = append(detected, Opportunity{
				Type:   PointExplosion,
				Actor:  impact.UserID,
				Target: signedPoints(impact.Delta),
			})
		}
	}
	return detected
}

func forecastKey(forecast forecasts.Forecast) string {
	return forecast.UserID + "|" + forecast.MatchID
}

func signedPoints(points int) string {
	return fmt.Sprintf("%+d", points)
}
