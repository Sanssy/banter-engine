package opportunities

import (
	"reflect"
	"testing"

	"github.com/Sanssy/banter-engine/internal/matches"
)

func TestDetectLiveUpdates(t *testing.T) {
	previous := []matches.Match{
		{
			MatchID:  "match-1",
			HomeTeam: "France",
			AwayTeam: "Espagne",
			Status:   "preMatch",
		},
	}
	current := []matches.Match{
		{
			MatchID:  "match-1",
			HomeTeam: "France",
			AwayTeam: "Espagne",
			Status:   "firstHalf",
			Score:    matches.Score{Home: 1},
			Events: []matches.Event{
				{ID: "event-1", Type: "goal", Time: "12'"},
				{ID: "event-2", Type: "substitution", Time: "20'"},
			},
		},
	}
	want := EnsureIdentities([]Opportunity{
		{Type: MatchStarted, Actor: "France - Espagne", MatchID: "match-1"},
		{Type: ScoreChanged, Actor: "France - Espagne", Target: "1-0", MatchID: "match-1"},
		{Type: GoalSwing, Actor: "France", Target: "France - Espagne (1-0)", MatchID: "match-1"},
		{Type: ImportantMatchEvent, Actor: "France - Espagne", Target: "goal 12'", MatchID: "match-1", EventID: "event-1"},
	})

	got := DetectLiveUpdates(previous, current)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DetectLiveUpdates() = %#v, want %#v", got, want)
	}
}

func TestDetectLiveUpdatesDetectsMatchEndWithoutRepeatingEvents(t *testing.T) {
	event := matches.Event{ID: "event-1", Type: "goal", Time: "90'"}
	previous := []matches.Match{{MatchID: "match-1", HomeTeam: "France", AwayTeam: "Espagne", Status: "secondHalf", Events: []matches.Event{event}}}
	current := []matches.Match{{MatchID: "match-1", HomeTeam: "France", AwayTeam: "Espagne", Status: "fullTime", Events: []matches.Event{event}}}

	want := EnsureIdentities([]Opportunity{{Type: MatchEnded, Actor: "France - Espagne", MatchID: "match-1"}})
	if got := DetectLiveUpdates(previous, current); !reflect.DeepEqual(got, want) {
		t.Fatalf("DetectLiveUpdates() = %#v, want %#v", got, want)
	}
}
