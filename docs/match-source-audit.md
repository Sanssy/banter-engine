# Match Source Audit

## Finding

Before iteration 031, match retrieval was not scoped to the configured challenge.

The runtime flow was:

```text
engine.runOnce
-> mpp.Client.GetMatches()
-> GET /championships-current-matches
-> all current MPP championships
-> opportunity detectors
```

`/championships-current-matches` is explicitly documented as returning current matches across championships. The configured `CHALLENGE_ID` was used for standings and forecasts, but it was not passed to match retrieval. Challenge scope was therefore lost before `Match` values were built.

## Previous Construction

`internal/mpp/client.go` decoded the global response into `map[string]*matchDTO`. Each non-null entry was converted directly to `matches.Match`, enriched with names from `/championship-clubs`, and returned to `engine.runOnce`. The engine then fetched forecasts and generated opportunities for every returned match.

This explains how club matches unrelated to the World Cup challenge could reach the opportunity engine.

## Challenge-Scoped Data Chain

The available MPP resources are sufficient to reconstruct the required matches:

```text
CHALLENGE_ID
-> GET /challenge/{challengeId}
-> challenge.gameSettings.championshipId
-> GET /championship-calendar/{championshipId}/nearest-game-weeks
-> nearestGameWeeks.currentGameWeek.matchesIds
-> POST /championship-match/summaries
-> map[matchId]MatchSummary
-> matches.Match
-> opportunity detectors
```

When `currentGameWeek` is omitted by an older API response, the client selects the previous or next game week whose date range contains the current time. Outside an active range it falls back to the next available game week, then the previous one.

## Implemented Invariant

`mpp.Client.GetMatches` now requires a challenge ID. It derives the championship and current match IDs through the chain above. When converting summaries, it iterates the calendar's `matchesIds`, not the keys returned by the summaries endpoint.

This final filtering step is deliberate: even if MPP returns an extra summary, a match absent from the challenge game week cannot enter the returned slice and therefore cannot reach opportunity detection.

`/championship-clubs` remains a global lookup used only to resolve team display names. It cannot introduce a match because it supplies no match IDs.

## Regression Coverage

`internal/mpp/client_test.go` verifies that:

- challenge `challenge-1` resolves championship `8`;
- the current game week is selected over previous and next game weeks;
- only the current game week's IDs are posted to `/championship-match/summaries`;
- a returned AS Saint-Etienne/Guingamp summary not requested by the calendar is ignored;
- the resulting match contains Canada and Qatar;
- an in-progress dated game week is selected when `currentGameWeek` is absent.

## Runtime Validation

Live dry-run validation requires the Raspberry Pi environment's `MPP_TOKEN` and was not executed in the development workspace. On the deployed host, run:

```sh
sudo systemctl stop banter-engine
cd /opt/banter-engine
set -a
. /etc/banter-engine.env
set +a
DRY_RUN=true ./banter-engine dry-run
```

Generated match opportunities must reference teams from the active World Cup game week. Restart the service after validation:

```sh
sudo systemctl start banter-engine
```
