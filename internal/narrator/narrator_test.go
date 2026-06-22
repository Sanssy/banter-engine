package narrator

import (
	"path/filepath"
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
	def := catalog.OpportunityDefinition{ID: string(op.Type), Name: "La foule avait tort"}

	want := banter.GenerateWithDefinition(op, def)
	got := (DeterministicNarrator{}).Narrate(op, def, narrative.CrowdWrong)
	if got != want {
		t.Fatalf("Narrate() = %q, want unchanged deterministic message %q", got, want)
	}
}

func TestDigestSelectsFiveExamplesForSingleOpportunity(t *testing.T) {
	library, err := narrative.LoadLibrary(filepath.Join("..", "..", "resources", "narratives", "examples.json"))
	if err != nil {
		t.Fatalf("LoadLibrary() error = %v", err)
	}
	cat, err := catalog.LoadCatalog(filepath.Join("..", "..", "resources", "opportunities.json"))
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}
	nar := &OllamaNarrator{examples: library}
	ops := []opportunities.Opportunity{{Type: opportunities.EveryoneWasWrong}}

	selected := nar.selectDigestExamples(ops, cat)
	if len(selected) != narrative.MaxExamples {
		t.Fatalf("selectDigestExamples() returned %d examples, want %d", len(selected), narrative.MaxExamples)
	}
	for _, example := range selected {
		if example.Angle != narrative.CrowdWrong {
			t.Fatalf("selectDigestExamples() returned incompatible example: %#v", example)
		}
	}
}

func TestWriteFewShotExamplesCapsPromptAtFive(t *testing.T) {
	examples := make([]narrative.Example, 6)
	for index := range examples {
		examples[index] = narrative.Example{
			Category: "Crowd",
			Angle:    narrative.CrowdWrong,
			Facts:    "faits",
			Message:  "message",
		}
	}

	var prompt strings.Builder
	writeFewShotExamples(&prompt, examples)
	if got := strings.Count(prompt.String(), " - faits :"); got != narrative.MaxExamples {
		t.Fatalf("prompt contains %d examples, want maximum %d", got, narrative.MaxExamples)
	}
	if strings.Contains(prompt.String(), "Exemple 6") {
		t.Fatalf("prompt contains a sixth example:\n%s", prompt.String())
	}
}

func TestBuildLivePromptIncludesFactsAndNarrativeAngle(t *testing.T) {
	op := opportunities.Opportunity{
		Type:   opportunities.AgainstTheCrowd,
		Actor:  "LeDaveCoinCoin",
		Target: "Paraguay - Turquie",
	}
	def := catalog.OpportunityDefinition{
		ID:          string(op.Type),
		Name:        "Seul contre tous",
		Description: "Une minorité avait anticipé le résultat.",
	}

	examples := []narrative.Example{
		{
			Category: "Crowd",
			Angle:    narrative.MinorityVictory,
			Facts:    "Trois joueurs avaient choisi le Paraguay.",
			Message:  "Trois visionnaires, le reste demande le replay.",
		},
	}
	prompt := buildLivePrompt(op, def, narrative.MinorityVictory, examples)
	for _, expected := range []string{
		"Événement : Seul contre tous",
		"Acteur principal : LeDaveCoinCoin",
		"Cible / contexte : Paraguay - Turquie",
		"Angle narratif : MinorityVictory",
		"Une minorité avait raison",
		"Exemple 1 - faits : Trois joueurs avaient choisi le Paraguay.",
		"Exemple 1 - message : Trois visionnaires, le reste demande le replay.",
	} {
		if !strings.Contains(prompt, expected) {
			t.Errorf("prompt does not contain %q:\n%s", expected, prompt)
		}
	}
}
