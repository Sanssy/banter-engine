# Opportunity Reconciliation Report

This report compares every opportunity emitted by the active detectors with the runtime catalog as it existed before iteration 029.

| Detector source | Opportunity emitted | Initially present in catalog? |
| --- | --- | --- |
| `internal/opportunities/crowd.go` | `CrowdFavorite` | yes |
| `internal/opportunities/crowd.go` | `CrowdTrap` | yes |
| `internal/opportunities/crowd.go` | `PopularMistake` | yes |
| `internal/opportunities/heartbreak.go` | `90thMinuteHeartbreak` | no |
| `internal/opportunities/heartbreak.go` | `AddedTimeDisaster` | yes |
| `internal/opportunities/heartbreak.go` | `LastMinuteHero` | no |
| `internal/opportunities/heartbreak.go` | `RedCardDisaster` | yes |
| `internal/opportunities/heartbreak.go` | `VARVictim` | yes |
| `internal/opportunities/live.go` | `ImportantMatchEvent` | no |
| `internal/opportunities/live.go` | `MatchEnded` | no |
| `internal/opportunities/live.go` | `MatchStarted` | no |
| `internal/opportunities/live.go` | `ScoreChanged` | no |
| `internal/opportunities/massacre.go` | `EveryoneWasWrong` | yes |
| `internal/opportunities/massacre.go` | `PredictionMassacre` | yes |
| `internal/opportunities/podium.go` | `PodiumFight` | no |
| `internal/opportunities/points.go` | `BiggestLoser` | no |
| `internal/opportunities/points.go` | `BiggestWinner` | no |
| `internal/opportunities/points.go` | `PointExplosion` | no |
| `internal/opportunities/prophet.go` | `AgainstTheCrowd` | yes |
| `internal/opportunities/prophet.go` | `TheChosenOne` | yes |
| `internal/opportunities/ranking.go` | `EnteredTop3` | yes |
| `internal/opportunities/ranking.go` | `ExitedTop3` | yes |
| `internal/opportunities/ranking.go` | `LastPlaceLocked` | yes |
| `internal/opportunities/ranking.go` | `LeaderUnderPressure` | yes |
| `internal/opportunities/ranking.go` | `RankingOvertake` | yes |
| `internal/opportunities/score_events.go` | `EqualizerChaos` | yes |
| `internal/opportunities/score_events.go` | `GoalSwing` | yes |
| `internal/opportunities/score_events.go` | `MatchTurnaround` | yes |
| `internal/opportunities/streaks.go` | `ColdStreak` | no |
| `internal/opportunities/streaks.go` | `HotStreak` | no |
| `internal/opportunities/surprises.go` | `HugeUpset` | no |
| `internal/opportunities/trends.go` | `ComebackSeason` | yes |
| `internal/opportunities/trends.go` | `FreeFall` | yes |
| `internal/opportunities/trends.go` | `RunawayLeader` | yes |
| `internal/rivalries/tracker.go` | `Dominance` | yes |
| `internal/rivalries/tracker.go` | `Nemesis` | yes |
| `internal/rivalries/tracker.go` | `Revenge` | yes |

## Result

- Active emitted opportunities: 37
- Initially present: 24
- Initially missing: 13
- Missing definitions added during iteration 029: 13
- Final catalog size: 113

The opportunity registry is now checked against the catalog both in tests and during engine startup.
