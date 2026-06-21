package opportunities

import "github.com/Sanssy/banter-engine/internal/forecasts"

func forecastUserName(forecast forecasts.Forecast) string {
	if forecast.UserName != "" {
		return forecast.UserName
	}
	return forecast.UserID
}
