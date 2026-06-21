package contextbuilder

import (
	"strings"
	"testing"

	"github.com/DSanoussy/banter-engine/internal/forecasts"
	"github.com/DSanoussy/banter-engine/internal/matches"
	"github.com/DSanoussy/banter-engine/internal/model"
	"github.com/DSanoussy/banter-engine/internal/opportunities"
)

func TestBuildCreatesStructuredContext(t *testing.T) {
	standings := []model.Standing{{UserID: "user-1", Name: "Julien", Rank: 1, Points: 100}}
	forecastData := []forecasts.Forecast{{UserID: "user-1", MatchID: "match-1", Points: 5}}
	matchData := []matches.Match{{MatchID: "match-1", HomeTeam: "France", AwayTeam: "Espagne"}}
	detected := []opportunities.Opportunity{{Type: opportunities.TheChosenOne, Actor: "Julien"}}

	ctx := Build(standings, forecastData, matchData, detected)
	standings[0].Name = "changed"
	if ctx.Standings[0].Name != "Julien" {
		t.Fatalf("Build() retained caller slice, got name %q", ctx.Standings[0].Name)
	}

	summary, err := ctx.Summary()
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	for _, value := range []string{"Julien", "match-1", "France", opportunities.TheChosenOne} {
		if !strings.Contains(summary, value) {
			t.Fatalf("Summary() = %q, missing %q", summary, value)
		}
	}
}
