package opportunities

import (
	"reflect"
	"testing"

	"github.com/Sanssy/banter-engine/internal/forecasts"
	"github.com/Sanssy/banter-engine/internal/matches"
)

func TestDetectSurprises(t *testing.T) {
	previous := matches.Match{MatchID: "match-1", Status: "secondHalf"}
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
	previous.HomeTeam = match.HomeTeam
	previous.AwayTeam = match.AwayTeam
	matchForecasts := []forecasts.Forecast{
		{UserID: "user-1", Prediction: matches.Score{Home: 0, Away: 1}},
		{UserID: "user-2", Prediction: matches.Score{Home: 2, Away: 0}},
	}
	want := EnsureIdentities([]Opportunity{
		{Type: HugeUpset, Actor: "Espagne", Target: "France", MatchID: "match-1"},
		{Type: EveryoneWasWrong, Actor: "France - Espagne", MatchID: "match-1"},
		{Type: PredictionMassacre, Actor: "France - Espagne", MatchID: "match-1"},
		{Type: CrowdFavorite, Actor: "France", Target: "90%", MatchID: "match-1"},
		{Type: PopularMistake, Actor: "France", Target: "France - Espagne", MatchID: "match-1"},
		{Type: TheChosenOne, Actor: "user-1", Target: "France - Espagne", MatchID: "match-1"},
		{Type: AgainstTheCrowd, Actor: "user-1", Target: "France - Espagne", MatchID: "match-1"},
	})

	got := DetectSurprises(previous, match, matchForecasts)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DetectSurprises() = %#v, want %#v", got, want)
	}
}

func TestDetectSurprisesIgnoresUnfinishedMatch(t *testing.T) {
	match := matches.Match{
		Status:          "preMatch",
		PredictionStats: matches.PredictionStats{Home: 0.90, Away: 0.10},
	}

	if got := DetectSurprises(match, match, nil); len(got) != 0 {
		t.Fatalf("DetectSurprises() returned %d opportunities, want 0", len(got))
	}
}

func TestDetectSurprisesIgnoresMissingPredictionStats(t *testing.T) {
	match := matches.Match{Status: "fullTime"}

	if got := DetectSurprises(matches.Match{}, match, nil); len(got) != 0 {
		t.Fatalf("DetectSurprises() returned %d opportunities, want 0", len(got))
	}
}

func TestDetectSurprisesDoesNotRepeatFinishedMatchOpportunities(t *testing.T) {
	match := matches.Match{
		MatchID:         "match-1",
		HomeTeam:        "France",
		AwayTeam:        "Espagne",
		Status:          "fullTime",
		Score:           matches.Score{Away: 1},
		Quotations:      matches.Quotations{Home: 120, Away: 300},
		PredictionStats: matches.PredictionStats{Home: 0.90, Away: 0.10},
	}
	if got := DetectSurprises(match, match, nil); len(got) != 0 {
		t.Fatalf("DetectSurprises() returned %d repeated opportunities, want 0", len(got))
	}
}
