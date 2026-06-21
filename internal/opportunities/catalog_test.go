package opportunities

import "testing"

func TestCatalogRegistersEveryOpportunityOnce(t *testing.T) {
	expected := []string{
		RankingOvertake,
		EnteredTop3,
		ExitedTop3,
		LeaderUnderPressure,
		LastPlaceLocked,
		ComebackSeason,
		FreeFall,
		RunawayLeader,
		PodiumFight,
		HugeUpset,
		EveryoneWasWrong,
		TheChosenOne,
		AgainstTheCrowd,
		PredictionMassacre,
		HotStreak,
		ColdStreak,
		MatchStarted,
		MatchEnded,
		ScoreChanged,
		ImportantMatchEvent,
		GoalSwing,
		MatchTurnaround,
		EqualizerChaos,
		BiggestWinner,
		BiggestLoser,
		PointExplosion,
		NinetiethMinuteHeartbreak,
		AddedTimeDisaster,
		LastMinuteHero,
		VARVictim,
		RedCardDisaster,
		Nemesis,
		Revenge,
		Dominance,
	}

	definitions := Definitions()
	if len(definitions) != len(expected) {
		t.Fatalf("Definitions() returned %d entries, want %d", len(definitions), len(expected))
	}

	seen := make(map[string]struct{}, len(definitions))
	for _, definition := range definitions {
		if _, duplicate := seen[definition.ID]; duplicate {
			t.Fatalf("duplicate opportunity ID %q", definition.ID)
		}
		seen[definition.ID] = struct{}{}
		if definition.Category == "" {
			t.Fatalf("opportunity %q has no category", definition.ID)
		}
		if definition.Description == "" {
			t.Fatalf("opportunity %q has no description", definition.ID)
		}
		metadata := definition.Detection
		if !metadata.RequiresStandings && !metadata.RequiresMatches && !metadata.RequiresForecasts && !metadata.RequiresLiveEvents {
			t.Fatalf("opportunity %q has no detection metadata", definition.ID)
		}
	}

	for _, id := range expected {
		if _, registered := seen[id]; !registered {
			t.Fatalf("opportunity %q is not registered", id)
		}
		if definition, found := FindDefinition(id); !found || definition.ID != id {
			t.Fatalf("FindDefinition(%q) = %#v, %t", id, definition, found)
		}
	}
}

func TestDefinitionsReturnsCopy(t *testing.T) {
	definitions := Definitions()
	definitions[0].ID = "changed"

	if Definitions()[0].ID == "changed" {
		t.Fatal("Definitions() exposes mutable catalog state")
	}
}
