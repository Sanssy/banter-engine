package opportunities

import (
	"reflect"
	"testing"

	"github.com/Sanssy/banter-engine/internal/model"
)

func TestDetectReturnsAllSupportedOpportunities(t *testing.T) {
	previous := []model.Standing{
		{UserID: "amine", Name: "Amine", Rank: 1, Points: 100},
		{UserID: "killian", Name: "Killian", Rank: 2, Points: 80},
		{UserID: "william", Name: "William", Rank: 3, Points: 70},
		{UserID: "sanssy", Name: "Sanssy", Rank: 4, Points: 65},
		{UserID: "julien", Name: "Julien", Rank: 5, Points: 50},
	}
	current := []model.Standing{
		{UserID: "amine", Name: "Amine", Rank: 1, Points: 90},
		{UserID: "sanssy", Name: "Sanssy", Rank: 2, Points: 88},
		{UserID: "killian", Name: "Killian", Rank: 3, Points: 82},
		{UserID: "william", Name: "William", Rank: 4, Points: 75},
		{UserID: "julien", Name: "Julien", Rank: 5, Points: 55},
	}
	want := []Opportunity{
		{Type: RankingOvertake, Actor: "Sanssy", Target: "Killian"},
		{Type: RankingOvertake, Actor: "Sanssy", Target: "William"},
		{Type: EnteredTop3, Actor: "Sanssy"},
		{Type: ExitedTop3, Actor: "William"},
		{Type: LeaderUnderPressure, Actor: "Amine", Target: "Sanssy"},
	}

	got := Detect(previous, current)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Detect() = %#v, want %#v", got, want)
	}
}

func TestDetectLastPlaceLockedOnlyWhenLastPlaceChanges(t *testing.T) {
	previous := []model.Standing{{UserID: "alice", Rank: 1}, {UserID: "bob", Rank: 2}}
	current := []model.Standing{{UserID: "bob", Rank: 1}, {UserID: "alice", Name: "Alice", Rank: 2}}

	got, ok := detectLastPlaceLocked(previous, current)
	if !ok || got != (Opportunity{Type: LastPlaceLocked, Actor: "Alice"}) {
		t.Fatalf("detectLastPlaceLocked() = %#v, %t", got, ok)
	}
	if _, ok := detectLastPlaceLocked(current, current); ok {
		t.Fatal("detectLastPlaceLocked() repeated unchanged last place")
	}
}

func TestDetectReturnsNoOpportunitiesWithoutPreviousStandings(t *testing.T) {
	current := []model.Standing{{UserID: "amine", Name: "Amine", Rank: 1, Points: 100}}

	if got := Detect(nil, current); len(got) != 0 {
		t.Fatalf("Detect() returned %d opportunities, want 0", len(got))
	}
}

func TestDetectRankingOvertakes(t *testing.T) {
	previous := []model.Standing{
		{UserID: "amine", Name: "Amine", Rank: 1},
		{UserID: "killian", Name: "Killian", Rank: 2},
		{UserID: "william", Name: "William", Rank: 3},
	}
	current := []model.Standing{
		{UserID: "william", Name: "William", Rank: 1},
		{UserID: "amine", Name: "Amine", Rank: 2},
		{UserID: "killian", Name: "Killian", Rank: 3},
	}
	want := []Opportunity{
		{Type: RankingOvertake, Actor: "William", Target: "Amine"},
		{Type: RankingOvertake, Actor: "William", Target: "Killian"},
	}

	got := detectRankingOvertakes(previous, current)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("detectRankingOvertakes() = %#v, want %#v", got, want)
	}
}

func TestDetectRankingOvertakesIgnoresUnchangedOrderAndNewUsers(t *testing.T) {
	previous := []model.Standing{
		{UserID: "amine", Name: "Amine", Rank: 1},
		{UserID: "killian", Name: "Killian", Rank: 2},
	}
	current := []model.Standing{
		{UserID: "amine", Name: "Amine", Rank: 1},
		{UserID: "killian", Name: "Killian", Rank: 2},
		{UserID: "william", Name: "William", Rank: 3},
	}

	if got := detectRankingOvertakes(previous, current); len(got) != 0 {
		t.Fatalf("detectRankingOvertakes() returned %d opportunities, want 0", len(got))
	}
}
