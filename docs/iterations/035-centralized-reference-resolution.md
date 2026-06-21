# Iteration 035 - Centralized Reference Resolution

## Goal

Prevent technical MPP club and user identifiers from appearing in user-facing opportunities.

## Implementation

- Resolve club identifiers from `GET /championship-clubs`.
- Resolve user identifiers from `GET /challenge-standings/users-standings`.
- Store resolved usernames on forecasts before they reach detectors.
- Resolve opportunity actors and targets once more before message generation as a defensive boundary.
- Log club lookups, user lookups, and resolved match construction.

## Validation

```bash
go test ./...
go vet ./...
git diff --check
```

Generated opportunity messages must not contain `mpp_championship_club_` or `user_` identifiers.
