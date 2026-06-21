package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DSanoussy/banter-engine/internal/model"
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
