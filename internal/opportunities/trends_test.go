package opportunities

import (
	"reflect"
	"testing"

	"github.com/DSanoussy/banter-engine/internal/model"
)

func TestDetectRankingTrends(t *testing.T) {
	previous := []model.Standing{
		{UserID: "leader", Name: "Amine", Rank: 1, Points: 100},
		{UserID: "second", Name: "Killian", Rank: 2, Points: 95},
		{UserID: "fall", Name: "William", Rank: 3, Points: 90},
		{UserID: "middle", Name: "Julien", Rank: 4, Points: 85},
		{UserID: "climb", Name: "Sanssy", Rank: 5, Points: 80},
		{UserID: "last", Name: "Pierre", Rank: 6, Points: 75},
	}
	current := []model.Standing{
		{UserID: "leader", Name: "Amine", Rank: 1, Points: 120},
		{UserID: "climb", Name: "Sanssy", Rank: 2, Points: 100},
		{UserID: "second", Name: "Killian", Rank: 3, Points: 98},
		{UserID: "middle", Name: "Julien", Rank: 4, Points: 90},
		{UserID: "last", Name: "Pierre", Rank: 5, Points: 80},
		{UserID: "fall", Name: "William", Rank: 6, Points: 78},
	}
	want := []Opportunity{
		{Type: ComebackSeason, Actor: "Sanssy", Target: "3"},
		{Type: FreeFall, Actor: "William", Target: "3"},
		{Type: RunawayLeader, Actor: "Amine", Target: "Sanssy"},
	}

	if got := DetectRankingTrends(previous, current); !reflect.DeepEqual(got, want) {
		t.Fatalf("DetectRankingTrends() = %#v, want %#v", got, want)
	}
}

func TestDetectRankingTrendsIgnoresSmallChanges(t *testing.T) {
	previous := []model.Standing{
		{UserID: "a", Name: "Amine", Rank: 1, Points: 100},
		{UserID: "b", Name: "Benoit", Rank: 2, Points: 95},
	}
	current := []model.Standing{
		{UserID: "a", Name: "Amine", Rank: 1, Points: 109},
		{UserID: "b", Name: "Benoit", Rank: 2, Points: 95},
	}

	if got := DetectRankingTrends(previous, current); len(got) != 0 {
		t.Fatalf("DetectRankingTrends() returned %d opportunities, want 0", len(got))
	}
}
