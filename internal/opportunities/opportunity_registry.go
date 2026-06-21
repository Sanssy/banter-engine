package opportunities

type Opportunity struct {
	Type   string
	Actor  string
	Target string
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

var registeredTypes = []string{
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
	return append([]string(nil), registeredTypes...)
}
