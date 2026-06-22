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

// dedupKey builds a composite key that keeps distinct actors/targets apart for the same type.
// (Type, Actor) when Actor is set, (Type, Target) when only Target is set, Type alone otherwise.
func dedupKey(op opportunities.Opportunity) string {
	if op.Actor != "" {
		return string(op.Type) + "\x00" + op.Actor
	}
	if op.Target != "" {
		return string(op.Type) + "\x00" + op.Target
	}
	return string(op.Type)
}

// SelectTop returns at most n opportunities, deduplicated by (Type, Actor/Target), ordered by catalog severity descending.
func SelectTop(ops []opportunities.Opportunity, cat *catalog.Catalog, n int) []opportunities.Opportunity {
	scored := make([]scoredOpportunity, 0, len(ops))
	for _, op := range ops {
		op = opportunities.EnsureIdentity(op)
		def, found := cat.FindByID(string(op.Type))
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
		key := dedupKey(s.op)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, s.op)
		if len(result) >= n {
			break
		}
	}
	return result
}
