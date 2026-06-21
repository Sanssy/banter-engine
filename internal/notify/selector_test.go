package notify

import (
	"testing"

	"github.com/Sanssy/banter-engine/internal/opportunities"
)

func TestDedupKey_DistinctActors(t *testing.T) {
	// Three HugeUpset with distinct actors must produce distinct keys.
	ops := []opportunities.Opportunity{
		{Type: opportunities.HugeUpset, Actor: "Turquie"},
		{Type: opportunities.HugeUpset, Actor: "Équateur"},
		{Type: opportunities.HugeUpset, Actor: "Tchéquie"},
	}
	keys := make(map[string]bool, len(ops))
	for _, op := range ops {
		k := dedupKey(op)
		if keys[k] {
			t.Fatalf("duplicate dedup key %q for actor %q", k, op.Actor)
		}
		keys[k] = true
	}
}

func TestDedupKey_SameActorSameType(t *testing.T) {
	// Same type AND same actor must share the same key (genuine duplicate).
	a := opportunities.Opportunity{Type: opportunities.HugeUpset, Actor: "Turquie"}
	b := opportunities.Opportunity{Type: opportunities.HugeUpset, Actor: "Turquie"}
	if dedupKey(a) != dedupKey(b) {
		t.Fatal("same type+actor should produce the same dedup key")
	}
}

func TestDedupKey_NoActorFallsBackToType(t *testing.T) {
	op := opportunities.Opportunity{Type: opportunities.PredictionMassacre}
	if dedupKey(op) != opportunities.PredictionMassacre {
		t.Fatal("empty actor and target should fall back to type as dedup key")
	}
}

func TestDedupKey_TargetUsedWhenNoActor(t *testing.T) {
	a := opportunities.Opportunity{Type: opportunities.Nemesis, Target: "Alice"}
	b := opportunities.Opportunity{Type: opportunities.Nemesis, Target: "Bob"}
	if dedupKey(a) == dedupKey(b) {
		t.Fatal("distinct targets with no actor should produce distinct dedup keys")
	}
}
