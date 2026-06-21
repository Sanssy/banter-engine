package opportunities

import (
	"reflect"
	"testing"

	"github.com/Sanssy/banter-engine/internal/forecasts"
	"github.com/Sanssy/banter-engine/internal/matches"
)

func TestDetectProphetsDetectsChosenOneAndAgainstTheCrowd(t *testing.T) {
	match := rareAwayWin()
	matchForecasts := []forecasts.Forecast{
		{UserID: "user-1", Prediction: matches.Score{Away: 1}},
		{UserID: "user-2", Prediction: matches.Score{Home: 1}},
	}
	want := []Opportunity{
		{Type: TheChosenOne, Actor: "user-1", Target: "France - Espagne"},
		{Type: AgainstTheCrowd, Actor: "user-1", Target: "France - Espagne"},
	}

	if got := detectProphets(match, matchForecasts); !reflect.DeepEqual(got, want) {
		t.Fatalf("detectProphets() = %#v, want %#v", got, want)
	}
}

func TestDetectProphetsEmitsAgainstTheCrowdForEachCorrectUser(t *testing.T) {
	match := rareAwayWin()
	matchForecasts := []forecasts.Forecast{
		{UserID: "user-1", Prediction: matches.Score{Away: 1}},
		{UserID: "user-2", Prediction: matches.Score{Home: 1, Away: 2}},
	}
	want := []Opportunity{
		{Type: AgainstTheCrowd, Actor: "user-1", Target: "France - Espagne"},
		{Type: AgainstTheCrowd, Actor: "user-2", Target: "France - Espagne"},
	}

	if got := detectProphets(match, matchForecasts); !reflect.DeepEqual(got, want) {
		t.Fatalf("detectProphets() = %#v, want %#v", got, want)
	}
}

func TestDetectProphetsUsesStrictFivePercentThreshold(t *testing.T) {
	match := rareAwayWin()
	match.PredictionStats = matches.PredictionStats{Home: 0.95, Away: 0.05}

	if got := detectProphets(match, []forecasts.Forecast{{UserID: "user-1", Prediction: matches.Score{Away: 1}}}); len(got) != 0 {
		t.Fatalf("detectProphets() returned %d opportunities at 5%%, want 0", len(got))
	}
}

func rareAwayWin() matches.Match {
	return matches.Match{
		HomeTeam:        "France",
		AwayTeam:        "Espagne",
		Score:           matches.Score{Away: 1},
		Status:          "fullTime",
		PredictionStats: matches.PredictionStats{Home: 0.96, Away: 0.04},
	}
}
