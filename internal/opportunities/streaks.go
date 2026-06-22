package opportunities

import (
	"sort"
	"strconv"

	"github.com/Sanssy/banter-engine/internal/forecasts"
)

const minimumStreakLength = 5

func DetectStreaks(previous, current []forecasts.Forecast) []Opportunity {
	previousStreaks := streaksByUser(previous)
	currentStreaks := streaksByUser(current)

	userIDs := make([]string, 0, len(currentStreaks))
	for userID := range currentStreaks {
		userIDs = append(userIDs, userID)
	}
	sort.Strings(userIDs)

	var detected []Opportunity
	for _, userID := range userIDs {
		currentStreak := currentStreaks[userID]
		previousStreak := previousStreaks[userID]
		if currentStreak.length < minimumStreakLength ||
			(currentStreak.kind == previousStreak.kind && previousStreak.length >= minimumStreakLength) {
			continue
		}
		detected = append(detected, Opportunity{
			Type:   currentStreak.kind,
			Actor:  currentStreak.userName,
			Target: strconv.Itoa(currentStreak.length),
		})
	}
	return detected
}

type streak struct {
	kind     string
	length   int
	userName string
}

func streaksByUser(history []forecasts.Forecast) map[string]streak {
	byUser := make(map[string][]forecasts.Forecast)
	for _, forecast := range history {
		byUser[forecast.UserID] = append(byUser[forecast.UserID], forecast)
	}

	result := make(map[string]streak, len(byUser))
	for userID, userForecasts := range byUser {
		sort.SliceStable(userForecasts, func(i, j int) bool {
			return userForecasts[i].MatchDate.Before(userForecasts[j].MatchDate)
		})
		streakType, length := activeStreak(userForecasts)
		result[userID] = streak{
			kind:     streakType,
			length:   length,
			userName: forecastUserName(userForecasts[len(userForecasts)-1]),
		}
	}
	return result
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
