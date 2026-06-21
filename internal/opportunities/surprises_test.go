package opportunities

import (
	"reflect"
	"testing"

	"github.com/DSanoussy/banter-engine/internal/forecasts"
	"github.com/DSanoussy/banter-engine/internal/matches"
)

func TestDetectSurprises(t *testing.T) {
	match := matches.Match{
		MatchID:  "match-1",
		HomeTeam: "France",
		AwayTeam: "Espagne",
		Score:    matches.Score{Home: 0, Away: 1},
		Status:   "fullTime",
		Quotations: matches.Quotations{
			Home: 120,
			Draw: 320,
			Away: 300,
		},
		PredictionStats: matches.PredictionStats{
			Home: 0.90,
			Draw: 0.06,
			Away: 0.04,
		},
	}
	matchForecasts := []forecasts.Forecast{
		{UserID: "user-1", Prediction: matches.Score{Home: 0, Away: 1}},
		{UserID: "user-2", Prediction: matches.Score{Home: 2, Away: 0}},
	}
	want := []Opportunity{
		{Type: HugeUpset, Actor: "Espagne", Target: "France"},
		{Type: EveryoneWasWrong, Actor: "France - Espagne"},
		{Type: TheChosenOne, Actor: "user-1", Target: "France - Espagne"},
		{Type: PredictionMassacre, Actor: "France - Espagne"},
	}

	got := DetectSurprises(match, matchForecasts)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DetectSurprises() = %#v, want %#v", got, want)
	}
}

func TestDetectSurprisesIgnoresUnfinishedMatch(t *testing.T) {
	match := matches.Match{
		Status:          "preMatch",
		PredictionStats: matches.PredictionStats{Home: 0.90, Away: 0.10},
	}

	if got := DetectSurprises(match, nil); len(got) != 0 {
		t.Fatalf("DetectSurprises() returned %d opportunities, want 0", len(got))
	}
}

func TestDetectSurprisesIgnoresMissingPredictionStats(t *testing.T) {
	match := matches.Match{Status: "fullTime"}

	if got := DetectSurprises(match, nil); len(got) != 0 {
		t.Fatalf("DetectSurprises() returned %d opportunities, want 0", len(got))
	}
}
