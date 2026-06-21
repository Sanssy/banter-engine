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
	case opportunities.EnteredTop3:
		return fmt.Sprintf("🏆 %s entre dans le top 3.", op.Actor)
	case opportunities.ExitedTop3:
		return fmt.Sprintf("📉 %s quitte le top 3.", op.Actor)
	case opportunities.LeaderUnderPressure:
		return fmt.Sprintf("🔥 %s voit %s revenir dans son rétroviseur.", op.Actor, op.Target)
	case opportunities.LastPlaceLocked:
		return fmt.Sprintf("🔒 %s conserve solidement la dernière place.", op.Actor)
	default:
		return fmt.Sprintf("%s: %s -> %s", op.Type, op.Actor, op.Target)
	}
}
