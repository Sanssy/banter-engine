package opportunities

import "strconv"

type OpportunityType string

type Opportunity struct {
	Type    OpportunityType
	Actor   string
	Target  string
	MatchID string
	EventID string
	Key     string
}

func EnsureIdentity(op Opportunity) Opportunity {
	if op.Key == "" {
		op.Key = identityPart(string(op.Type)) + identityPart(op.MatchID) +
			identityPart(op.EventID) + identityPart(op.Actor) + identityPart(op.Target)
	}
	return op
}

func EnsureIdentities(ops []Opportunity) []Opportunity {
	for i := range ops {
		ops[i] = EnsureIdentity(ops[i])
	}
	return ops
}

func identityPart(value string) string {
	return strconv.Itoa(len(value)) + ":" + value
}

const (
	RankingOvertake           = "RankingOvertake"
	EnteredTop3               = "EnteredTop3"
	ExitedTop3                = "ExitedTop3"
	LeaderUnderPressure       = "LeaderUnderPressure"
	LastPlaceLocked           = "LastPlaceLocked"
	ComebackSeason            = "ComebackSeason"
	FreeFall                  = "FreeFall"
	RunawayLeader             = "RunawayLeader"
	PodiumFight               = "PodiumFight"
	HugeUpset                 = "HugeUpset"
	EveryoneWasWrong          = "EveryoneWasWrong"
	TheChosenOne              = "TheChosenOne"
	AgainstTheCrowd           = "AgainstTheCrowd"
	CrowdFavorite             = "CrowdFavorite"
	CrowdTrap                 = "CrowdTrap"
	PopularMistake            = "PopularMistake"
	PredictionMassacre        = "PredictionMassacre"
	HotStreak                 = "HotStreak"
	ColdStreak                = "ColdStreak"
	MatchStarted              = "MatchStarted"
	MatchEnded                = "MatchEnded"
	ScoreChanged              = "ScoreChanged"
	ImportantMatchEvent       = "ImportantMatchEvent"
	GoalSwing                 = "GoalSwing"
	MatchTurnaround           = "MatchTurnaround"
	EqualizerChaos            = "EqualizerChaos"
	BiggestWinner             = "BiggestWinner"
	BiggestLoser              = "BiggestLoser"
	PointExplosion            = "PointExplosion"
	NinetiethMinuteHeartbreak = "90thMinuteHeartbreak"
	AddedTimeDisaster         = "AddedTimeDisaster"
	LastMinuteHero            = "LastMinuteHero"
	VARVictim                 = "VARVictim"
	RedCardDisaster           = "RedCardDisaster"
	Nemesis                   = "Nemesis"
	Revenge                   = "Revenge"
	Dominance                 = "Dominance"
)

var registeredTypes = []OpportunityType{
	RankingOvertake,
	EnteredTop3,
	ExitedTop3,
	LeaderUnderPressure,
	LastPlaceLocked,
	ComebackSeason,
	FreeFall,
	RunawayLeader,
	PodiumFight,
	HugeUpset,
	EveryoneWasWrong,
	TheChosenOne,
	AgainstTheCrowd,
	CrowdFavorite,
	CrowdTrap,
	PopularMistake,
	PredictionMassacre,
	HotStreak,
	ColdStreak,
	MatchStarted,
	MatchEnded,
	ScoreChanged,
	ImportantMatchEvent,
	GoalSwing,
	MatchTurnaround,
	EqualizerChaos,
	BiggestWinner,
	BiggestLoser,
	PointExplosion,
	NinetiethMinuteHeartbreak,
	AddedTimeDisaster,
	LastMinuteHero,
	VARVictim,
	RedCardDisaster,
	Nemesis,
	Revenge,
	Dominance,
}

func RegisteredTypes() []string {
	result := make([]string, len(registeredTypes))
	for i, opportunityType := range registeredTypes {
		result[i] = string(opportunityType)
	}
	return result
}
