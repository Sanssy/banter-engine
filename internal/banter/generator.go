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
	case opportunities.HugeUpset:
		return fmt.Sprintf("⚠️ %s fait tomber le favori %s.", op.Actor, op.Target)
	case opportunities.EveryoneWasWrong:
		return fmt.Sprintf("📉 %s : l'intelligence collective a pris un coup.", op.Actor)
	case opportunities.TheChosenOne:
		return fmt.Sprintf("🔮 %s était le seul à y croire sur %s.", op.Actor, op.Target)
	case opportunities.AgainstTheCrowd:
		return fmt.Sprintf("🎯 %s a défié la foule et avait raison sur %s.", op.Actor, op.Target)
	case opportunities.PredictionMassacre:
		return fmt.Sprintf("☠️ Extinction des pronostics détectée sur %s.", op.Actor)
	case opportunities.HotStreak:
		return fmt.Sprintf("🔥 %s enchaîne %s pronostics réussis.", op.Actor, op.Target)
	case opportunities.ColdStreak:
		return fmt.Sprintf("🥶 %s reste sur %s échecs consécutifs.", op.Actor, op.Target)
	case opportunities.MatchStarted:
		return fmt.Sprintf("⚽ Coup d'envoi pour %s.", op.Actor)
	case opportunities.MatchEnded:
		return fmt.Sprintf("🏁 Fin du match %s.", op.Actor)
	case opportunities.ScoreChanged:
		return fmt.Sprintf("⚽ %s : le score passe à %s.", op.Actor, op.Target)
	case opportunities.ImportantMatchEvent:
		return fmt.Sprintf("🚨 %s : %s.", op.Actor, op.Target)
	case opportunities.NinetiethMinuteHeartbreak:
		return fmt.Sprintf("💔 %s perd son prono parfait dans les derniers instants de %s.", op.Actor, op.Target)
	case opportunities.VARVictim:
		return fmt.Sprintf("📺 La VAR brise le pronostic de %s sur %s.", op.Actor, op.Target)
	case opportunities.RedCardDisaster:
		return fmt.Sprintf("🟥 Le rouge met le pronostic de %s en danger sur %s.", op.Actor, op.Target)
	case opportunities.Nemesis:
		return fmt.Sprintf("👹 %s devient la Némésis officielle de %s.", op.Actor, op.Target)
	case opportunities.Revenge:
		return fmt.Sprintf("🗡️ %s prend sa revanche sur %s.", op.Actor, op.Target)
	case opportunities.Dominance:
		return fmt.Sprintf("👑 %s domine désormais sa rivalité avec %s.", op.Actor, op.Target)
	default:
		return fmt.Sprintf("%s: %s -> %s", op.Type, op.Actor, op.Target)
	}
}
