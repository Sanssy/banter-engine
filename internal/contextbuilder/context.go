package contextbuilder

import (
	"encoding/json"
	"fmt"

	"github.com/DSanoussy/banter-engine/internal/forecasts"
	"github.com/DSanoussy/banter-engine/internal/matches"
	"github.com/DSanoussy/banter-engine/internal/model"
	"github.com/DSanoussy/banter-engine/internal/opportunities"
)

type Context struct {
	Standings     []model.Standing            `json:"standings"`
	Forecasts     []forecasts.Forecast        `json:"forecasts"`
	Matches       []matches.Match             `json:"matches"`
	Opportunities []opportunities.Opportunity `json:"opportunities"`
}

func Build(
	standings []model.Standing,
	forecastData []forecasts.Forecast,
	matchData []matches.Match,
	detected []opportunities.Opportunity,
) Context {
	return Context{
		Standings:     append([]model.Standing(nil), standings...),
		Forecasts:     append([]forecasts.Forecast(nil), forecastData...),
		Matches:       append([]matches.Match(nil), matchData...),
		Opportunities: append([]opportunities.Opportunity(nil), detected...),
	}
}

func (c Context) Summary() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("encode banter context: %w", err)
	}
	return string(data), nil
}
