MPP Context

Purpose

This document describes the information required from MPP in order to detect Banter Opportunities.

The objective is not to document the API.

The objective is to understand which business data Banter Engine depends on.

⸻

Core Concepts Needed From MPP

Players

We need information about league participants.

Examples:

* player identifier
* player display name

⸻

Standings

We need the current league standings.

Examples:

* rank
* total points
* points gap

Used by:

* LeaderUnderPressure
* LostFirstPlace
* RankingOvertake
* LastPlaceLocked
* PodiumFight

⸻

Predictions

We need predictions submitted by players.

Examples:

* predicted winner
* predicted score
* prediction timestamp

Used by:

* PerfectPrediction
* AlmostPerfectPrediction
* PredictionDestroyed
* TheChosenOne
* AgainstTheCrowd

⸻

Points

We need the points awarded for each prediction.

Examples:

* points earned
* points lost
* total points

Used by:

* WinningStreak
* LosingStreak
* RapidClimb
* FreeFall

⸻

Matches

We need information about matches associated with the competition.

Examples:

* match identifier
* teams
* status
* competition
* kickoff time

Used by nearly every Banter Opportunity.

## Historical Data

We need access to historical information.

Examples:

- previous standings
- previous points
- previous predictions

Used by:

- WinningStreak
- LosingStreak
- RapidClimb
- FreeFall
- BottomKing

---

## Prediction Distribution

We need to understand how predictions are distributed across the league.

Examples:

- percentage predicting a home win
- percentage predicting a draw
- percentage predicting an away win
- players choosing minority predictions

Used by:

- EveryoneWasWrong
- TheChosenOne
- AgainstTheCrowd
- PredictionMassacre

---

## Match Events

Some opportunities require live match events.

Examples:

- goals
- red cards
- penalties
- VAR decisions
- disallowed goals

Used by:

- PredictionDestroyed
- LateGoalChaos
- 90thMinuteHeartbreak
- PredictionSaved

---

## External Data Sources

MPP is expected to provide league-specific data such as standings, predictions, points and users.

However, some opportunities may require external data sources.

Examples:

- live score API
- match event API
- odds API

These sources may be needed for live banter, upset detection and match event analysis.
