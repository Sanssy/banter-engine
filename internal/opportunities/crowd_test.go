package opportunities

import (
	"reflect"
	"testing"

	"github.com/Sanssy/banter-engine/internal/matches"
)

func TestDetectCrowdIntelligenceDetectsPopularMistake(t *testing.T) {
	match := matches.Match{
		HomeTeam: "France",
		AwayTeam: "Espagne",
		Score:    matches.Score{Away: 1},
		Status:   "fullTime",
		Quotations: matches.Quotations{
			Home: 120,
			Away: 300,
		},
		PredictionStats: matches.PredictionStats{Home: 0.70, Draw: 0.10, Away: 0.20},
	}
	want := []Opportunity{
		{Type: CrowdFavorite, Actor: "France", Target: "70%"},
		{Type: PopularMistake, Actor: "France", Target: "France - Espagne"},
	}

	if got := detectCrowdIntelligence(match); !reflect.DeepEqual(got, want) {
		t.Fatalf("detectCrowdIntelligence() = %#v, want %#v", got, want)
	}
}

func TestDetectCrowdIntelligenceDetectsCrowdTrap(t *testing.T) {
	match := matches.Match{
		HomeTeam: "France",
		AwayTeam: "Espagne",
		Quotations: matches.Quotations{
			Home: 120,
			Away: 300,
		},
		PredictionStats: matches.PredictionStats{Home: 0.20, Draw: 0.15, Away: 0.65},
	}
	want := []Opportunity{
		{Type: CrowdFavorite, Actor: "Espagne", Target: "65%"},
		{Type: CrowdTrap, Actor: "Espagne", Target: "France"},
	}

	if got := detectCrowdIntelligence(match); !reflect.DeepEqual(got, want) {
		t.Fatalf("detectCrowdIntelligence() = %#v, want %#v", got, want)
	}
}

func TestDetectCrowdIntelligenceRequiresSixtyPercent(t *testing.T) {
	match := matches.Match{PredictionStats: matches.PredictionStats{Home: 0.59, Draw: 0.21, Away: 0.20}}
	if got := detectCrowdIntelligence(match); len(got) != 0 {
		t.Fatalf("detectCrowdIntelligence() returned %d opportunities, want 0", len(got))
	}
}

func TestDetectCrowdTransitionsOnlyWhenThresholdIsCrossed(t *testing.T) {
	previous := matches.Match{HomeTeam: "France", PredictionStats: matches.PredictionStats{Home: 0.59, Away: 0.41}}
	current := matches.Match{HomeTeam: "France", PredictionStats: matches.PredictionStats{Home: 0.60, Away: 0.40}}

	want := []Opportunity{{Type: CrowdFavorite, Actor: "France", Target: "60%"}}
	if got := detectCrowdTransitions(previous, current); !reflect.DeepEqual(got, want) {
		t.Fatalf("detectCrowdTransitions() = %#v, want %#v", got, want)
	}
	if got := detectCrowdTransitions(current, current); len(got) != 0 {
		t.Fatalf("detectCrowdTransitions() repeated %d opportunities, want 0", len(got))
	}
}
