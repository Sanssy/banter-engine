package rivalries

import (
	"reflect"
	"testing"

	"github.com/Sanssy/banter-engine/internal/model"
	"github.com/Sanssy/banter-engine/internal/opportunities"
)

func TestUpdateTracksHeadToHeadWithoutCountingUnchangedStandings(t *testing.T) {
	standings := []model.Standing{
		{UserID: "a", Name: "Amine", Rank: 1, Points: 10},
		{UserID: "b", Name: "Benoit", Rank: 2, Points: 5},
	}

	state, detected := Update(standings, nil)
	if len(detected) != 0 {
		t.Fatalf("initial Update() returned %d opportunities, want 0", len(detected))
	}
	if state[0].WinsA != 1 || state[0].WinsB != 0 {
		t.Fatalf("initial score = %d-%d, want 1-0", state[0].WinsA, state[0].WinsB)
	}

	unchangedState, detected := Update(standings, state)
	if len(detected) != 0 {
		t.Fatalf("unchanged Update() returned %d opportunities, want 0", len(detected))
	}
	if !reflect.DeepEqual(unchangedState, state) {
		t.Fatalf("unchanged Update() changed state: %#v", unchangedState)
	}
}

func TestUpdateDetectsRevengeNemesisAndDominance(t *testing.T) {
	first := []model.Standing{
		{UserID: "a", Name: "Amine", Rank: 1, Points: 10},
		{UserID: "b", Name: "Benoit", Rank: 2, Points: 5},
	}
	state, _ := Update(first, nil)

	revengeStandings := []model.Standing{
		{UserID: "b", Name: "Benoit", Rank: 1, Points: 12},
		{UserID: "a", Name: "Amine", Rank: 2, Points: 10},
	}
	state, detected := Update(revengeStandings, state)
	wantRevenge := opportunities.EnsureIdentities([]opportunities.Opportunity{
		{Type: opportunities.Revenge, Actor: "Benoit", Target: "Amine"},
	})
	if !reflect.DeepEqual(detected, wantRevenge) {
		t.Fatalf("revenge Update() = %#v, want %#v", detected, wantRevenge)
	}

	for _, points := range []int{13, 14} {
		revengeStandings[0].Points = points
		state, detected = Update(revengeStandings, state)
	}
	wantMilestones := opportunities.EnsureIdentities([]opportunities.Opportunity{
		{Type: opportunities.Nemesis, Actor: "Benoit", Target: "Amine"},
	})
	if !reflect.DeepEqual(detected, wantMilestones) {
		t.Fatalf("milestone Update() = %#v, want %#v", detected, wantMilestones)
	}

	revengeStandings[0].Points = 15
	_, detected = Update(revengeStandings, state)
	wantDominance := opportunities.EnsureIdentities([]opportunities.Opportunity{
		{Type: opportunities.Dominance, Actor: "Benoit", Target: "Amine"},
	})
	if !reflect.DeepEqual(detected, wantDominance) {
		t.Fatalf("dominance Update() = %#v, want %#v", detected, wantDominance)
	}
}
