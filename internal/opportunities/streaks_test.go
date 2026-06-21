package opportunities

import (
	"reflect"
	"testing"
	"time"

	"github.com/DSanoussy/banter-engine/internal/forecasts"
)

func TestDetectStreaks(t *testing.T) {
	start := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	history := []forecasts.Forecast{
		{UserID: "hot-user", MatchDate: start.Add(1 * time.Hour), Points: 3},
		{UserID: "cold-user", MatchDate: start.Add(5 * time.Hour), Points: 0},
		{UserID: "hot-user", MatchDate: start.Add(3 * time.Hour), Points: 5},
		{UserID: "cold-user", MatchDate: start.Add(1 * time.Hour), Points: 2},
		{UserID: "hot-user", MatchDate: start.Add(2 * time.Hour), Points: 3},
		{UserID: "cold-user", MatchDate: start.Add(2 * time.Hour), Points: 0},
		{UserID: "hot-user", MatchDate: start.Add(5 * time.Hour), Points: 2},
		{UserID: "cold-user", MatchDate: start.Add(3 * time.Hour), Points: 0},
		{UserID: "hot-user", MatchDate: start.Add(4 * time.Hour), Points: 3},
		{UserID: "cold-user", MatchDate: start.Add(4 * time.Hour), Points: 0},
		{UserID: "cold-user", MatchDate: start.Add(6 * time.Hour), Points: 0},
	}
	want := []Opportunity{
		{Type: ColdStreak, Actor: "cold-user", Target: "5"},
		{Type: HotStreak, Actor: "hot-user", Target: "5"},
	}

	got := DetectStreaks(history)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DetectStreaks() = %#v, want %#v", got, want)
	}
}

func TestDetectStreaksRequiresFiveConsecutiveResults(t *testing.T) {
	history := []forecasts.Forecast{
		{UserID: "user-1", Points: 3},
		{UserID: "user-1", Points: 3},
		{UserID: "user-1", Points: 0},
		{UserID: "user-1", Points: 3},
		{UserID: "user-1", Points: 3},
		{UserID: "user-1", Points: 3},
		{UserID: "user-1", Points: 3},
	}

	if got := DetectStreaks(history); len(got) != 0 {
		t.Fatalf("DetectStreaks() returned %d opportunities, want 0", len(got))
	}
}
