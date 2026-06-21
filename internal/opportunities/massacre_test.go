package opportunities

import (
	"reflect"
	"testing"

	"github.com/DSanoussy/banter-engine/internal/matches"
)

func TestDetectMassFailures(t *testing.T) {
	match := matches.Match{
		HomeTeam: "France",
		AwayTeam: "Espagne",
		Score:    matches.Score{Away: 1},
		Status:   "fullTime",
		PredictionStats: matches.PredictionStats{
			Home: 0.81,
			Draw: 0.09,
			Away: 0.10,
		},
	}
	want := []Opportunity{
		{Type: EveryoneWasWrong, Actor: "France - Espagne"},
		{Type: PredictionMassacre, Actor: "France - Espagne"},
	}

	if got := detectMassFailures(match); !reflect.DeepEqual(got, want) {
		t.Fatalf("detectMassFailures() = %#v, want %#v", got, want)
	}
}

func TestDetectMassFailuresUsesStrictThresholds(t *testing.T) {
	match := matches.Match{
		Score:  matches.Score{Home: 1},
		Status: "fullTime",
		PredictionStats: matches.PredictionStats{
			Home: 0.50,
			Draw: 0.30,
			Away: 0.20,
		},
	}

	if got := detectMassFailures(match); len(got) != 0 {
		t.Fatalf("detectMassFailures() returned %d opportunities at the boundary, want 0", len(got))
	}

	match.Score = matches.Score{Away: 1}
	match.PredictionStats = matches.PredictionStats{Home: 0.80, Away: 0.20}
	got := detectMassFailures(match)
	want := []Opportunity{{Type: PredictionMassacre, Actor: " - "}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("detectMassFailures() = %#v, want %#v", got, want)
	}
}
