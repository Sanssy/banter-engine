package opportunities

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
	NinetiethMinuteHeartbreak = "90thMinuteHeartbreak"
	VARVictim                 = "VARVictim"
	RedCardDisaster           = "RedCardDisaster"
	Nemesis                   = "Nemesis"
	Revenge                   = "Revenge"
	Dominance                 = "Dominance"
)

type Category string

const (
	CategoryRanking    Category = "ranking"
	CategoryPrediction Category = "prediction"
	CategoryStreak     Category = "streak"
	CategoryLive       Category = "live"
	CategoryRivalry    Category = "rivalry"
)

type DetectionMetadata struct {
	RequiresStandings        bool
	RequiresPreviousSnapshot bool
	RequiresMatches          bool
	RequiresForecasts        bool
	RequiresLiveEvents       bool
}

type Definition struct {
	ID          string
	Category    Category
	Description string
	Detection   DetectionMetadata
}

var catalog = []Definition{
	{RankingOvertake, CategoryRanking, "A player moves ahead of another player.", standingsHistory()},
	{EnteredTop3, CategoryRanking, "A player enters the top three.", standingsHistory()},
	{ExitedTop3, CategoryRanking, "A player exits the top three.", standingsHistory()},
	{LeaderUnderPressure, CategoryRanking, "The leader's points advantage shrinks.", standingsHistory()},
	{LastPlaceLocked, CategoryRanking, "The same player remains in last place.", standingsHistory()},
	{ComebackSeason, CategoryRanking, "A player climbs at least three positions.", standingsHistory()},
	{FreeFall, CategoryRanking, "A player drops at least three positions.", standingsHistory()},
	{RunawayLeader, CategoryRanking, "The leader materially increases the lead.", standingsHistory()},
	{PodiumFight, CategoryRanking, "Fourth place closes in on the podium.", standingsHistory()},
	{HugeUpset, CategoryPrediction, "The match favorite fails to win.", matchDetection(false)},
	{EveryoneWasWrong, CategoryPrediction, "More than 80 percent backed the same wrong outcome.", matchDetection(false)},
	{TheChosenOne, CategoryPrediction, "Exactly one player picked a rare correct outcome.", matchDetection(true)},
	{AgainstTheCrowd, CategoryPrediction, "A player picked a rare correct outcome.", matchDetection(true)},
	{PredictionMassacre, CategoryPrediction, "More than half of predictions are wrong.", matchDetection(false)},
	{HotStreak, CategoryStreak, "A player has at least five consecutive successes.", forecastHistory()},
	{ColdStreak, CategoryStreak, "A player has at least five consecutive failures.", forecastHistory()},
	{MatchStarted, CategoryLive, "A match transitions to live.", liveDetection(false)},
	{MatchEnded, CategoryLive, "A live match reaches full time.", liveDetection(false)},
	{ScoreChanged, CategoryLive, "A live match score changes.", liveDetection(false)},
	{ImportantMatchEvent, CategoryLive, "A new important match event occurs.", liveDetection(false)},
	{GoalSwing, CategoryLive, "A goal changes the live score.", liveDetection(false)},
	{MatchTurnaround, CategoryLive, "The winning side changes during a match.", liveDetection(false)},
	{EqualizerChaos, CategoryLive, "A goal brings a live match level.", liveDetection(false)},
	{NinetiethMinuteHeartbreak, CategoryLive, "A late goal destroys an exact forecast.", liveDetection(true)},
	{VARVictim, CategoryLive, "A VAR decision destroys a forecast.", liveDetection(true)},
	{RedCardDisaster, CategoryLive, "A red card threatens a forecast.", liveDetection(true)},
	{Nemesis, CategoryRivalry, "A player records a third head-to-head win.", standingsHistory()},
	{Revenge, CategoryRivalry, "A player reverses the latest head-to-head result.", standingsHistory()},
	{Dominance, CategoryRivalry, "A player leads a rivalry by at least three wins.", standingsHistory()},
}

func Definitions() []Definition {
	return append([]Definition(nil), catalog...)
}

func FindDefinition(id string) (Definition, bool) {
	for _, definition := range catalog {
		if definition.ID == id {
			return definition, true
		}
	}
	return Definition{}, false
}

func standingsHistory() DetectionMetadata {
	return DetectionMetadata{RequiresStandings: true, RequiresPreviousSnapshot: true}
}

func matchDetection(requiresForecasts bool) DetectionMetadata {
	return DetectionMetadata{RequiresMatches: true, RequiresForecasts: requiresForecasts}
}

func forecastHistory() DetectionMetadata {
	return DetectionMetadata{RequiresForecasts: true, RequiresPreviousSnapshot: true}
}

func liveDetection(requiresForecasts bool) DetectionMetadata {
	return DetectionMetadata{
		RequiresPreviousSnapshot: true,
		RequiresMatches:          true,
		RequiresForecasts:        requiresForecasts,
		RequiresLiveEvents:       true,
	}
}
