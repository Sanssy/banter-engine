# Reference Resolution Audit

## Club flow

```text
GET /championship-clubs
  -> clubDTO.lang[locale].name
  -> references.Resolver
  -> matches.Match.HomeTeam / AwayTeam
  -> detectors
  -> opportunity message
```

The club reference was decoded, but the DTO only read root-level `name` and `shortName` fields. The MPP payload also stores localized names under `lang`, so unresolved clubs fell back to their technical identifiers. The client now registers localized names centrally before building matches.

## User flow

```text
GET /challenge-standings/users-standings
  -> model.Standing.UserID / Name
  -> references.Resolver
  -> forecasts.Forecast.UserName
  -> detectors
  -> opportunity message
```

Forecasts previously carried only `UserID`, even though standings already contained the username. User-based detectors therefore used the technical identifier as their actor. Forecasts now carry the resolved username and detectors use it for user-facing fields.

## Output boundary

Immediately before banter generation, opportunity actors and targets pass through the same resolver. This prevents a technical MPP identifier from leaking if a detector receives an older snapshot or an unresolved model.

Unknown technical club and user identifiers are rendered as `Equipe inconnue` and `Participant inconnu`. The lookup log records `found=false` so missing reference data remains diagnosable without exposing the identifier to users.
