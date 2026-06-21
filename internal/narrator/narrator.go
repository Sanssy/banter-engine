package narrator

import (
	"github.com/Sanssy/banter-engine/internal/catalog"
	"github.com/Sanssy/banter-engine/internal/opportunities"
)

// Narrator reformulates a single live opportunity into a short, natural message.
type Narrator interface {
	Narrate(op opportunities.Opportunity, def catalog.OpportunityDefinition) string
}

// DigestNarrator summarizes a batch of overnight opportunities into a morning digest.
type DigestNarrator interface {
	Summarize(ops []opportunities.Opportunity, cat *catalog.Catalog) string
}
