package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DSanoussy/banter-engine/internal/forecasts"
	"github.com/DSanoussy/banter-engine/internal/matches"
	"github.com/DSanoussy/banter-engine/internal/model"
	"github.com/DSanoussy/banter-engine/internal/rivalries"
)

func SaveStandings(path string, standings []model.Standing) error {
	data, err := json.MarshalIndent(standings, "", "  ")
	if err != nil {
		return fmt.Errorf("encode standings: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create snapshot directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write standings snapshot: %w", err)
	}

	return nil
}

func LoadStandings(path string) ([]model.Standing, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []model.Standing{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read standings snapshot: %w", err)
	}

	var standings []model.Standing
	if err := json.Unmarshal(data, &standings); err != nil {
		return nil, fmt.Errorf("decode standings snapshot: %w", err)
	}

	return standings, nil
}

func SaveMatches(path string, matches []matches.Match) error {
	data, err := json.MarshalIndent(matches, "", "  ")
	if err != nil {
		return fmt.Errorf("encode matches: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create snapshot directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write matches snapshot: %w", err)
	}
	return nil
}

func LoadMatches(path string) ([]matches.Match, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []matches.Match{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read matches snapshot: %w", err)
	}

	var result []matches.Match
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode matches snapshot: %w", err)
	}
	return result, nil
}

func SaveRivalries(path string, state []rivalries.Rivalry) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encode rivalries: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create snapshot directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write rivalries snapshot: %w", err)
	}
	return nil
}

func LoadRivalries(path string) ([]rivalries.Rivalry, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []rivalries.Rivalry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read rivalries snapshot: %w", err)
	}

	var state []rivalries.Rivalry
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("decode rivalries snapshot: %w", err)
	}
	return state, nil
}

func SaveForecasts(path string, forecastData []forecasts.Forecast) error {
	data, err := json.MarshalIndent(forecastData, "", "  ")
	if err != nil {
		return fmt.Errorf("encode forecasts: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create snapshot directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write forecasts snapshot: %w", err)
	}
	return nil
}

func LoadForecasts(path string) ([]forecasts.Forecast, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []forecasts.Forecast{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read forecasts snapshot: %w", err)
	}

	var forecastData []forecasts.Forecast
	if err := json.Unmarshal(data, &forecastData); err != nil {
		return nil, fmt.Errorf("decode forecasts snapshot: %w", err)
	}
	return forecastData, nil
}
