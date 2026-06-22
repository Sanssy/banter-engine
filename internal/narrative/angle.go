package narrative

import "github.com/Sanssy/banter-engine/internal/opportunities"

type Angle string

const (
	CrowdWrong      Angle = "CrowdWrong"
	MinorityVictory Angle = "MinorityVictory"
	FallFromGrace   Angle = "FallFromGrace"
	Rise            Angle = "Rise"
	Curse           Angle = "Curse"
	Dominance       Angle = "Dominance"
)

func ForOpportunity(opportunityType string) Angle {
	switch opportunityType {
	case opportunities.EveryoneWasWrong,
		opportunities.PredictionMassacre,
		opportunities.PopularMistake,
		opportunities.CrowdTrap:
		return CrowdWrong

	case opportunities.TheChosenOne,
		opportunities.AgainstTheCrowd:
		return MinorityVictory

	case opportunities.HugeUpset,
		opportunities.ExitedTop3,
		opportunities.LastPlaceLocked,
		opportunities.FreeFall,
		opportunities.BiggestLoser,
		opportunities.NinetiethMinuteHeartbreak,
		opportunities.AddedTimeDisaster,
		opportunities.VARVictim,
		opportunities.RedCardDisaster:
		return FallFromGrace

	case opportunities.RankingOvertake,
		opportunities.EnteredTop3,
		opportunities.ComebackSeason,
		opportunities.PodiumFight,
		opportunities.BiggestWinner,
		opportunities.PointExplosion,
		opportunities.LastMinuteHero,
		opportunities.Revenge,
		opportunities.MatchStarted,
		opportunities.ScoreChanged,
		opportunities.ImportantMatchEvent,
		opportunities.GoalSwing,
		opportunities.MatchTurnaround,
		opportunities.EqualizerChaos:
		return Rise

	case opportunities.ColdStreak,
		opportunities.Nemesis:
		return Curse

	case opportunities.LeaderUnderPressure,
		opportunities.RunawayLeader,
		opportunities.CrowdFavorite,
		opportunities.HotStreak,
		opportunities.Dominance,
		opportunities.MatchEnded:
		return Dominance
	default:
		return ""
	}
}

func (a Angle) Guidance() string {
	switch a {
	case CrowdWrong:
		return "La majorité s'est trompée : souligne l'échec collectif."
	case MinorityVictory:
		return "Une minorité avait raison : mets en valeur ceux qui ont osé se démarquer."
	case FallFromGrace:
		return "Un favori ou un joueur chute : raconte le renversement sans cruauté."
	case Rise:
		return "Une progression ou un retournement se produit : souligne l'élan positif."
	case Curse:
		return "Une série noire se poursuit : souligne la répétition avec légèreté."
	case Dominance:
		return "Une domination ou une prise de pouvoir se confirme."
	default:
		return ""
	}
}
