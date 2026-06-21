package opportunities

import (
	"reflect"
	"testing"

	"github.com/Sanssy/banter-engine/internal/model"
)

func TestDetectPodiumFight(t *testing.T) {
	previous := []model.Standing{
		{UserID: "third", Name: "William", Rank: 3, Points: 80},
		{UserID: "fourth", Name: "Sanssy", Rank: 4, Points: 70},
	}
	current := []model.Standing{
		{UserID: "third", Name: "William", Rank: 3, Points: 85},
		{UserID: "fourth", Name: "Sanssy", Rank: 4, Points: 82},
	}
	want := []Opportunity{{Type: PodiumFight, Actor: "Sanssy", Target: "William"}}

	if got := detectPodiumFight(previous, current); !reflect.DeepEqual(got, want) {
		t.Fatalf("detectPodiumFight() = %#v, want %#v", got, want)
	}
}

func TestDetectPodiumFightDoesNotRepeatUnchangedFight(t *testing.T) {
	previous := []model.Standing{
		{UserID: "third", Name: "William", Rank: 3, Points: 85},
		{UserID: "fourth", Name: "Sanssy", Rank: 4, Points: 82},
	}
	current := []model.Standing{
		{UserID: "third", Name: "William", Rank: 3, Points: 86},
		{UserID: "fourth", Name: "Sanssy", Rank: 4, Points: 83},
	}

	if got := detectPodiumFight(previous, current); len(got) != 0 {
		t.Fatalf("detectPodiumFight() returned %d opportunities, want 0", len(got))
	}
}
