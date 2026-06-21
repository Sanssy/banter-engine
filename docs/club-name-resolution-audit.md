# Club Name Resolution Audit

## Data Flow

Match summaries contain `home.clubId` and `away.clubId`. `mpp.Client.GetMatches` resolves those IDs through `/championship-clubs` before constructing `matches.Match`. Every detector then uses `Match.HomeTeam` and `Match.AwayTeam`; no detector reads a club ID directly.

## Failure Point

The previous DTO accepted only this response shape:

```json
{"championshipClubs":{"club-id":{"name":{"fr-FR":"France"}}}}
```

MPP can return the club map directly at the response root:

```json
{"club-id":{"name":{"fr-FR":"France"}}}
```

JSON decoding did not fail: the direct club keys were treated as unknown struct fields, leaving `ChampionshipClubs` empty. `clubName` consequently used its fallback argument, which is the technical club ID. That fallback propagated unchanged through the domain match, detectors, and banter generator.

## Resolution

`clubsResponse` now accepts both direct and wrapped maps. Name resolution prefers French, then English localized names, then `shortName`, and uses the technical ID only when the reference genuinely lacks a display name.

The MPP client regression test uses a direct-root club response and verifies that the resulting World Cup match contains `Canada` and `Qatar`. A user-facing opportunity test verifies that generated banter contains the resolved team name and never the `mpp_championship_club_` prefix.
