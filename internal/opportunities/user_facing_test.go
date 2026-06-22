package opportunities_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Sanssy/banter-engine/internal/banter"
	"github.com/Sanssy/banter-engine/internal/catalog"
	"github.com/Sanssy/banter-engine/internal/forecasts"
	"github.com/Sanssy/banter-engine/internal/matches"
	"github.com/Sanssy/banter-engine/internal/opportunities"
)

func TestGeneratedOpportunityUsesResolvedClubName(t *testing.T) {
	match := matches.Match{
		HomeTeam: "Argentine",
		AwayTeam: "France",
		Status:   "fullTime",
		Score:    matches.Score{Home: 1},
		Quotations: matches.Quotations{
			Home: 50,
			Away: 100,
		},
		PredictionStats: matches.PredictionStats{Home: 0.95, Draw: 0.03, Away: 0.02},
	}
	detected := opportunities.DetectSurprises(matches.Match{Status: "secondHalf"}, match, nil)
	if len(detected) == 0 {
		t.Fatal("DetectSurprises() returned no opportunity")
	}

	message := banter.GenerateWithDefinition(detected[0], catalog.OpportunityDefinition{ID: string(detected[0].Type)})
	if !strings.Contains(message, "Argentine") {
		t.Fatalf("generated message %q does not contain resolved club name", message)
	}
	if strings.Contains(message, "mpp_championship_club_") {
		t.Fatalf("generated message %q contains technical club ID", message)
	}
}

func TestGeneratedOpportunityUsesResolvedUserName(t *testing.T) {
	start := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	history := make([]forecasts.Forecast, 5)
	for i := range history {
		history[i] = forecasts.Forecast{
			UserID:    "user_11291094",
			UserName:  "LeDaveCoinCoin",
			MatchDate: start.Add(time.Duration(i) * time.Hour),
		}
	}

	detected := opportunities.DetectStreaks(nil, history)
	if len(detected) != 1 {
		t.Fatalf("DetectStreaks() returned %d opportunities, want 1", len(detected))
	}

	message := banter.GenerateWithDefinition(detected[0], catalog.OpportunityDefinition{ID: string(detected[0].Type)})
	if !strings.Contains(message, "LeDaveCoinCoin") {
		t.Fatalf("generated message %q does not contain resolved user name", message)
	}
	if strings.Contains(message, "user_") {
		t.Fatalf("generated message %q contains technical user ID", message)
	}
}
