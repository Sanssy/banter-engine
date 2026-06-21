package context

import (
	"strings"
	"testing"

	"github.com/DSanoussy/banter-engine/internal/catalog"
	"github.com/DSanoussy/banter-engine/internal/opportunities"
)

func TestBuildLLMContextIncludesDefinitionAndMetadata(t *testing.T) {
	opportunity := opportunities.Opportunity{
		Type:   opportunities.RankingOvertake,
		Actor:  "Sanssy",
		Target: "William",
	}
	definition := catalog.OpportunityDefinition{
		ID:          opportunities.RankingOvertake,
		Category:    "Ranking",
		Severity:    2,
		Description: "A player overtakes another player.",
		Tags:        []string{"ranking", "leaderboard"},
	}

	result := BuildLLMContext(opportunity, definition)
	for _, value := range []string{
		definition.ID,
		definition.Category,
		definition.Description,
		"severity",
		"ranking",
		opportunity.Actor,
		opportunity.Target,
	} {
		if !strings.Contains(result, value) {
			t.Fatalf("BuildLLMContext() = %q, missing %q", result, value)
		}
	}
}
