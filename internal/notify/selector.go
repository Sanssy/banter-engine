package notify

import (
	"sort"

	"github.com/Sanssy/banter-engine/internal/catalog"
	"github.com/Sanssy/banter-engine/internal/opportunities"
)

// MaxNotificationsPerRun is the maximum number of messages sent in a single cycle.
const MaxNotificationsPerRun = 5

type scoredOpportunity struct {
	op    opportunities.Opportunity
	score int
}

// SelectTop returns at most n opportunities, deduplicated by type, ordered by catalog severity descending.
func SelectTop(ops []opportunities.Opportunity, cat *catalog.Catalog, n int) []opportunities.Opportunity {
	scored := make([]scoredOpportunity, 0, len(ops))
	for _, op := range ops {
		def, found := cat.FindByID(op.Type)
		if !found {
			continue
		}
		scored = append(scored, scoredOpportunity{op, def.Severity})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	seen := make(map[string]bool, n)
	result := make([]opportunities.Opportunity, 0, n)
	for _, s := range scored {
		if seen[s.op.Type] {
			continue
		}
		seen[s.op.Type] = true
		result = append(result, s.op)
		if len(result) >= n {
			break
		}
	}
	return result
}
