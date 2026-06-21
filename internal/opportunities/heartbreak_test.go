package opportunities

import (
	"reflect"
	"testing"

	"github.com/Sanssy/banter-engine/internal/forecasts"
	"github.com/Sanssy/banter-engine/internal/matches"
)

func TestDetectHeartbreaksDetectsLateGoal(t *testing.T) {
	previous := matches.Match{MatchID: "match-1", Score: matches.Score{Home: 1}}
	current := matches.Match{
		MatchID:  "match-1",
		HomeTeam: "France",
		AwayTeam: "Espagne",
		Score:    matches.Score{Home: 1, Away: 1},
		Events:   []matches.Event{{ID: "goal-1", Type: "goal", Detail: "goal", Time: "90' +3"}},
	}
	matchForecasts := []forecasts.Forecast{
		{UserID: "user-1", Prediction: matches.Score{Home: 1}},
	}
	want := []Opportunity{
		{Type: NinetiethMinuteHeartbreak, Actor: "user-1", Target: "France - Espagne"},
	}

	if got := DetectHeartbreaks(previous, current, matchForecasts); !reflect.DeepEqual(got, want) {
		t.Fatalf("DetectHeartbreaks() = %#v, want %#v", got, want)
	}
}

func TestDetectHeartbreaksDetectsVARVictim(t *testing.T) {
	previous := matches.Match{MatchID: "match-1", Score: matches.Score{Home: 1, Away: 1}}
	current := matches.Match{
		MatchID:  "match-1",
		HomeTeam: "France",
		AwayTeam: "Espagne",
		Score:    matches.Score{Home: 1},
		Events:   []matches.Event{{ID: "var-1", Type: "goal", Detail: "var", Time: "70'"}},
	}
	matchForecasts := []forecasts.Forecast{
		{UserID: "user-1", Prediction: matches.Score{Home: 1, Away: 1}},
	}
	want := []Opportunity{{Type: VARVictim, Actor: "user-1", Target: "France - Espagne"}}

	if got := DetectHeartbreaks(previous, current, matchForecasts); !reflect.DeepEqual(got, want) {
		t.Fatalf("DetectHeartbreaks() = %#v, want %#v", got, want)
	}
}

func TestDetectHeartbreaksDetectsRedCardDisaster(t *testing.T) {
	previous := matches.Match{MatchID: "match-1"}
	current := matches.Match{
		MatchID:  "match-1",
		HomeTeam: "France",
		AwayTeam: "Espagne",
		Events: []matches.Event{
			{ID: "card-1", Type: "booking", Detail: "straightRed", Side: "home", Time: "60'"},
		},
	}
	matchForecasts := []forecasts.Forecast{
		{UserID: "user-1", Prediction: matches.Score{Home: 2, Away: 0}},
		{UserID: "user-2", Prediction: matches.Score{Home: 0, Away: 1}},
	}
	want := []Opportunity{{Type: RedCardDisaster, Actor: "user-1", Target: "France - Espagne"}}

	if got := DetectHeartbreaks(previous, current, matchForecasts); !reflect.DeepEqual(got, want) {
		t.Fatalf("DetectHeartbreaks() = %#v, want %#v", got, want)
	}
}
