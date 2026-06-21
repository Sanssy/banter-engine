package opportunities

import (
	"reflect"
	"testing"

	"github.com/DSanoussy/banter-engine/internal/forecasts"
)

func TestCalculateAndDetectPointImpacts(t *testing.T) {
	previous := []forecasts.Forecast{
		{UserID: "user-1", MatchID: "match-1", Points: 0},
		{UserID: "user-2", MatchID: "match-1", Points: 5},
		{UserID: "user-3", MatchID: "match-1", Points: 1},
	}
	current := []forecasts.Forecast{
		{UserID: "user-1", MatchID: "match-1", Points: 7},
		{UserID: "user-2", MatchID: "match-1", Points: 0},
		{UserID: "user-3", MatchID: "match-1", Points: 3},
	}
	wantImpacts := []pointImpact{
		{UserID: "user-1", MatchID: "match-1", PreviousPoints: 0, CurrentPoints: 7, Delta: 7},
		{UserID: "user-2", MatchID: "match-1", PreviousPoints: 5, CurrentPoints: 0, Delta: -5},
		{UserID: "user-3", MatchID: "match-1", PreviousPoints: 1, CurrentPoints: 3, Delta: 2},
	}
	if got := calculatePointImpacts(previous, current); !reflect.DeepEqual(got, wantImpacts) {
		t.Fatalf("calculatePointImpacts() = %#v, want %#v", got, wantImpacts)
	}

	wantOpportunities := []Opportunity{
		{Type: BiggestWinner, Actor: "user-1", Target: "+7"},
		{Type: BiggestLoser, Actor: "user-2", Target: "-5"},
		{Type: PointExplosion, Actor: "user-1", Target: "+7"},
	}
	if got := DetectPointImpacts(previous, current); !reflect.DeepEqual(got, wantOpportunities) {
		t.Fatalf("DetectPointImpacts() = %#v, want %#v", got, wantOpportunities)
	}
}

func TestDetectPointImpactsIgnoresFirstObservationAndUnchangedPoints(t *testing.T) {
	previous := []forecasts.Forecast{{UserID: "user-1", MatchID: "match-1", Points: 3}}
	current := []forecasts.Forecast{
		{UserID: "user-1", MatchID: "match-1", Points: 3},
		{UserID: "user-2", MatchID: "match-1", Points: 5},
	}

	if got := DetectPointImpacts(previous, current); len(got) != 0 {
		t.Fatalf("DetectPointImpacts() returned %d opportunities, want 0", len(got))
	}
}
