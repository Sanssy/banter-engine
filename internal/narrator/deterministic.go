package narrator

import (
	"github.com/Sanssy/banter-engine/internal/banter"
	"github.com/Sanssy/banter-engine/internal/catalog"
	"github.com/Sanssy/banter-engine/internal/notify"
	"github.com/Sanssy/banter-engine/internal/opportunities"
)

// DeterministicNarrator uses existing template-based generation — always available.
type DeterministicNarrator struct{}

func (DeterministicNarrator) Narrate(op opportunities.Opportunity, def catalog.OpportunityDefinition) string {
	return banter.GenerateWithDefinition(op, def)
}

func (DeterministicNarrator) Summarize(ops []opportunities.Opportunity, cat *catalog.Catalog) string {
	return notify.NightSummary(ops, cat)
}
