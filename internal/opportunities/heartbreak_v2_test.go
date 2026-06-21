package opportunities

import (
	"reflect"
	"testing"

	"github.com/Sanssy/banter-engine/internal/forecasts"
	"github.com/Sanssy/banter-engine/internal/matches"
)

func TestDetectLatePointImpacts(t *testing.T) {
	previousMatch := matches.Match{MatchID: "match-1"}
	currentMatch := matches.Match{
		MatchID:  "match-1",
		HomeTeam: "France",
		AwayTeam: "Espagne",
		Events:   []matches.Event{{ID: "goal-1", Type: "goal", Time: "90' +3"}},
	}
	previousForecasts := []forecasts.Forecast{
		{UserID: "user-1", MatchID: "match-1", Points: 5},
		{UserID: "user-2", MatchID: "match-1", Points: 0},
	}
	currentForecasts := []forecasts.Forecast{
		{UserID: "user-1", MatchID: "match-1", Points: 0},
		{UserID: "user-2", MatchID: "match-1", Points: 5},
	}
	want := []Opportunity{
		{Type: AddedTimeDisaster, Actor: "user-1", Target: "France - Espagne"},
		{Type: LastMinuteHero, Actor: "user-2", Target: "France - Espagne"},
	}

	got := DetectLatePointImpacts(previousMatch, currentMatch, previousForecasts, currentForecasts)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DetectLatePointImpacts() = %#v, want %#v", got, want)
	}
}

func TestDetectLatePointImpactsIgnoresEventsBeforeNinetiethMinute(t *testing.T) {
	previousMatch := matches.Match{MatchID: "match-1"}
	currentMatch := matches.Match{
		MatchID: "match-1",
		Events:  []matches.Event{{ID: "goal-1", Type: "goal", Time: "89'"}},
	}
	previousForecasts := []forecasts.Forecast{{UserID: "user-1", MatchID: "match-1", Points: 0}}
	currentForecasts := []forecasts.Forecast{{UserID: "user-1", MatchID: "match-1", Points: 5}}

	if got := DetectLatePointImpacts(previousMatch, currentMatch, previousForecasts, currentForecasts); len(got) != 0 {
		t.Fatalf("DetectLatePointImpacts() returned %d opportunities, want 0", len(got))
	}
}
