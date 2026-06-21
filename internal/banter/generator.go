package banter

import (
	"fmt"

	"github.com/DSanoussy/banter-engine/internal/opportunities"
)

func Generate(op opportunities.Opportunity) string {
	switch op.Type {
	case opportunities.RankingOvertake:
		return fmt.Sprintf(
			"📈 %s dépasse %s.\n\n%s aperçoit désormais son dos.",
			op.Actor,
			op.Target,
			op.Target,
		)
	default:
		return fmt.Sprintf("%s: %s -> %s", op.Type, op.Actor, op.Target)
	}
}
