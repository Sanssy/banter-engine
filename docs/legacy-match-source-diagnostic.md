# Legacy Match Source Diagnostic

## Evidence

The captured API responses establish the following facts:

- challenge `mpp_challenge_UDKDDH27` uses championship `8`;
- current game week `2` contains match IDs `2608265` through `2608288`;
- `/championship-match/summaries` returns those same IDs with championship `8`;
- `/championships-current-matches` returns the observed `2569xxx` IDs and unrelated Ligue 1/Ligue 2 clubs.

Deleting `data/` does not change the result. Snapshots are outputs of the current run and are not the source of the legacy IDs.

## Code Path Audit

In the current source tree, `engine.runOnce` obtains its current match slice only from:

```text
mpp.Client.GetMatches(CHALLENGE_ID)
```

That method calls only:

```text
GET  /challenge/{challengeId}
GET  /championship-calendar/{championshipId}/nearest-game-weeks
POST /championship-match/summaries
GET  /championship-clubs
```

The string `/championships-current-matches` is absent from production Go code. No bootstrap file, mock, alternate provider, or forecast-to-match reconstruction exists. `LoadMatches` reads only the previous snapshot for comparisons and never replaces the newly fetched `matches` slice. `LoadForecasts` returns forecast history only.

## Injection Point

The `2569xxx` payload captured in `diagnostic.txt` is byte-for-byte characteristic of `/championships-current-matches`, the endpoint used by `GetMatches()` before commit `b1d99aa`.

Therefore, if runtime logs do not contain the new `component=mpp` challenge/calendar/summaries messages and still persist `2569xxx`, the running process is an older pre-031 executable or a different executable path. Removing snapshots cannot affect that executable, so it immediately calls the old global endpoint and recreates the same data.

There is one secondary boundary to verify: the current converter prefers `summary.matchId` over the requested map key. The captured World Cup summaries show matching `2608xxx` values, so this boundary does not explain the supplied evidence. The diagnostic logs print both values to confirm this in the running process.

## Runtime Decision Table

| Observed log sequence | Conclusion |
| --- | --- |
| No `startup ... revision=...` or `component=mpp` diagnostic lines | Old or different binary is running. |
| Calendar ID is `2608xxx`, `summary_match_id` becomes `2569xxx` | MPP summary conversion boundary injects the ID. |
| Summary and resolved IDs are `2608xxx`, detector input becomes `2569xxx` | An in-process mutation exists between the client and engine; current static audit found none. |
| Detector and persistence IDs are `2608xxx` but old messages still appear | The messages come from another running process or Discord webhook producer. |

## Correction To Apply After Confirmation

The evidence currently points to deployment drift. The corrective action would be to install the binary built from commit `b1d99aa` or later at the exact systemd `ExecStart` path, reload/restart the service, and verify the logged VCS revision and executable path.

If live logs instead prove that `summary.matchId` differs from the requested calendar key, the later functional fix should preserve the requested calendar ID and reject summaries whose championship or game week does not match the selected challenge scope.

Neither correction is applied by this diagnostic ticket.
