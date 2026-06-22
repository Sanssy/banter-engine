package narrator

import (
	"strings"
	"testing"

	"github.com/Sanssy/banter-engine/internal/banter"
	"github.com/Sanssy/banter-engine/internal/catalog"
	"github.com/Sanssy/banter-engine/internal/narrative"
	"github.com/Sanssy/banter-engine/internal/opportunities"
)

func TestDeterministicNarratorOutputDoesNotChangeWithAngle(t *testing.T) {
	op := opportunities.Opportunity{
		Type:   opportunities.EveryoneWasWrong,
		Actor:  "Paraguay - Turquie",
		Target: "92%",
	}
	def := catalog.OpportunityDefinition{ID: op.Type, Name: "La foule avait tort"}

	want := banter.GenerateWithDefinition(op, def)
	got := (DeterministicNarrator{}).Narrate(op, def, narrative.CrowdWrong)
	if got != want {
		t.Fatalf("Narrate() = %q, want unchanged deterministic message %q", got, want)
	}
}

func TestBuildLivePromptIncludesFactsAndNarrativeAngle(t *testing.T) {
	op := opportunities.Opportunity{
		Type:   opportunities.AgainstTheCrowd,
		Actor:  "LeDaveCoinCoin",
		Target: "Paraguay - Turquie",
	}
	def := catalog.OpportunityDefinition{
		ID:          op.Type,
		Name:        "Seul contre tous",
		Description: "Une minorité avait anticipé le résultat.",
	}

	prompt := buildLivePrompt(op, def, narrative.MinorityVictory)
	for _, expected := range []string{
		"Événement : Seul contre tous",
		"Acteur principal : LeDaveCoinCoin",
		"Cible / contexte : Paraguay - Turquie",
		"Angle narratif : MinorityVictory",
		"Une minorité avait raison",
	} {
		if !strings.Contains(prompt, expected) {
			t.Errorf("prompt does not contain %q:\n%s", expected, prompt)
		}
	}
}
