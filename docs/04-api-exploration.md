API Exploration

Purpose

This document maps Banter Engine domain needs to available data sources.

The objective is not to document every endpoint.

The objective is to determine whether Banter Engine can obtain the information required to detect Banter Opportunities.

⸻

Data Source Strategy

MPP API

MPP is considered the source of truth for:

* players
* leagues
* standings
* predictions
* points
* prediction history

MPP answers the question:

“How is the prediction league evolving?”

⸻

Live Match Data API

A dedicated football data provider will likely be required for:

* live scores
* goals
* red cards
* penalties
* substitutions
* VAR decisions
* match timeline

This source answers:

“What is currently happening on the pitch?”

⸻

Odds API (Optional)

Odds may be useful for detecting unexpected outcomes.

Examples:

* GiantKiller
* EveryoneWasWrong
* MiraclePrediction

This source answers:

“How surprising is the result?”

⸻

Domain Needs Mapping

Domain Need	Required Data	Expected Source	Status
Standings	Rankings, points, positions	MPP	Unknown
Players	User information	MPP	Unknown
Predictions	Match predictions	MPP	Unknown
Historical Rankings	Previous standings	MPP or local storage	Unknown
Historical Predictions	Previous predictions	MPP or local storage	Unknown
Prediction Distribution	All predictions for a match	MPP	Unknown
Match Results	Final scores	MPP or Live API	Unknown
Live Scores	Current score	Live API	Unknown
Match Events	Goals, cards, VAR	Live API	Unknown
Odds	Match odds	Odds API	Optional

⸻

Candidate Live Data Providers

API-Football

Potentially the strongest candidate.

Capabilities:

* live scores
* match events
* standings
* odds
* historical data

Live data is updated frequently and includes events such as goals and cards. (GitHub)

Status: Candidate

⸻

Other Providers

To be evaluated:

* Football Data
* Sofascore integrations
* Futmetrics
* Open football data providers

Status: Pending evaluation. (GitHub)

⸻

MPP API Exploration

Authentication

Unknown

Standings

Unknown endpoint

Status: To investigate

Predictions

Unknown endpoint

Status: To investigate

Points

Unknown endpoint

Status: To investigate

Matches

Unknown endpoint

Status: To investigate

⸻

Open Questions

1. Can MPP provide all league predictions for a given match?
2. Can MPP provide prediction history?
3. Can MPP provide ranking history?
4. Does MPP expose live scores?
5. Does MPP expose match events?
6. Can MPP alone power the MVP?

Initial Findings

The MPP API appears significantly richer than initially expected.

Confirmed capabilities:

* Challenge standings
* User rankings
* Match predictions
* Forecast points
* Match summaries
* Match details
* Match odds
* Prediction distribution
* Match event timeline

Potentially detectable opportunities:

* LeaderUnderPressure
* LostFirstPlace
* RankingOvertake
* LastPlaceLocked
* PerfectPrediction
* AlmostPerfectPrediction
* EveryoneWasWrong
* AgainstTheCrowd
* GiantKiller
* PredictionMassacre

Open question:

The exact content of eventsTimeline must be inspected to determine whether live opportunities such as:

* LateGoalChaos
* PredictionDestroyed
* 90thMinuteHeartbreak

can be detected using MPP alone.

MVP Feasibility Assessment

Current assessment: HIGH

The MPP API appears to provide enough information to build a first Banter Engine MVP without relying on external providers.

Available capabilities:

* League standings
* Historical standings
* User rankings
* Match schedules
* Match summaries
* Match odds
* Prediction distribution
* Forecast points

Potential MVP opportunities:

* LeaderUnderPressure
* LostFirstPlace
* RankingOvertake
* LastPlaceLocked
* PerfectPrediction
* AlmostPerfectPrediction
* EveryoneWasWrong
* AgainstTheCrowd
* GiantKiller
* PredictionMassacre

External live score providers may still be required for advanced opportunities involving:

* live goals
* red cards
* VAR events
* late match drama
* prediction reversals during a match

However, these are not required for the first MVP.
