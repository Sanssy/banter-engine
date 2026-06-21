package opportunities

import (
	"reflect"
	"testing"

	"github.com/DSanoussy/banter-engine/internal/model"
)

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

	got := DetectRankingOvertakes(previous, current)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DetectRankingOvertakes() = %#v, want %#v", got, want)
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

	if got := DetectRankingOvertakes(previous, current); len(got) != 0 {
		t.Fatalf("DetectRankingOvertakes() returned %d opportunities, want 0", len(got))
	}
}
