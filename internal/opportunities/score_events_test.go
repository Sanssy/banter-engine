package opportunities

import (
	"reflect"
	"testing"

	"github.com/DSanoussy/banter-engine/internal/matches"
)

func TestDetectScoreEventsDetectsEqualizer(t *testing.T) {
	previous := matches.Match{HomeTeam: "France", AwayTeam: "Espagne", Score: matches.Score{Home: 1}}
	current := matches.Match{HomeTeam: "France", AwayTeam: "Espagne", Score: matches.Score{Home: 1, Away: 1}}
	want := []Opportunity{
		{Type: GoalSwing, Actor: "Espagne", Target: "France - Espagne (1-1)"},
		{Type: EqualizerChaos, Actor: "Espagne", Target: "France - Espagne"},
	}

	if got := detectScoreEvents(previous, current); !reflect.DeepEqual(got, want) {
		t.Fatalf("detectScoreEvents() = %#v, want %#v", got, want)
	}
}

func TestDetectScoreEventsDetectsTurnaround(t *testing.T) {
	previous := matches.Match{HomeTeam: "France", AwayTeam: "Espagne", Score: matches.Score{Away: 1}}
	current := matches.Match{HomeTeam: "France", AwayTeam: "Espagne", Score: matches.Score{Home: 2, Away: 1}}
	want := []Opportunity{
		{Type: GoalSwing, Actor: "France", Target: "France - Espagne (2-1)"},
		{Type: MatchTurnaround, Actor: "France", Target: "Espagne"},
	}

	if got := detectScoreEvents(previous, current); !reflect.DeepEqual(got, want) {
		t.Fatalf("detectScoreEvents() = %#v, want %#v", got, want)
	}
}

func TestDetectScoreEventsIgnoresUnchangedScore(t *testing.T) {
	match := matches.Match{Score: matches.Score{Home: 1, Away: 1}}
	if got := detectScoreEvents(match, match); len(got) != 0 {
		t.Fatalf("detectScoreEvents() returned %d opportunities, want 0", len(got))
	}
}
