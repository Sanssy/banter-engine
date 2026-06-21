package opportunities_test

import (
	"strings"
	"testing"

	"github.com/Sanssy/banter-engine/internal/banter"
	"github.com/Sanssy/banter-engine/internal/catalog"
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
	detected := opportunities.DetectSurprises(match, nil)
	if len(detected) == 0 {
		t.Fatal("DetectSurprises() returned no opportunity")
	}

	message := banter.GenerateWithDefinition(detected[0], catalog.OpportunityDefinition{ID: detected[0].Type})
	if !strings.Contains(message, "Argentine") {
		t.Fatalf("generated message %q does not contain resolved club name", message)
	}
	if strings.Contains(message, "mpp_championship_club_") {
		t.Fatalf("generated message %q contains technical club ID", message)
	}
}
