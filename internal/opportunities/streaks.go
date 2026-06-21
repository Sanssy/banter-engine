package opportunities

import (
	"sort"
	"strconv"

	"github.com/Sanssy/banter-engine/internal/forecasts"
)

const minimumStreakLength = 5

func DetectStreaks(history []forecasts.Forecast) []Opportunity {
	byUser := make(map[string][]forecasts.Forecast)
	for _, forecast := range history {
		byUser[forecast.UserID] = append(byUser[forecast.UserID], forecast)
	}

	userIDs := make([]string, 0, len(byUser))
	for userID := range byUser {
		userIDs = append(userIDs, userID)
	}
	sort.Strings(userIDs)

	var detected []Opportunity
	for _, userID := range userIDs {
		forecasts := byUser[userID]
		sort.SliceStable(forecasts, func(i, j int) bool {
			return forecasts[i].MatchDate.Before(forecasts[j].MatchDate)
		})

		streakType, length := activeStreak(forecasts)
		if length >= minimumStreakLength {
			detected = append(detected, Opportunity{
				Type:   streakType,
				Actor:  userID,
				Target: strconv.Itoa(length),
			})
		}
	}

	return detected
}

func activeStreak(history []forecasts.Forecast) (string, int) {
	if len(history) == 0 {
		return "", 0
	}

	isSuccess := history[len(history)-1].Points > 0
	length := 0
	for i := len(history) - 1; i >= 0; i-- {
		if (history[i].Points > 0) != isSuccess {
			break
		}
		length++
	}

	if isSuccess {
		return HotStreak, length
	}
	return ColdStreak, length
}
