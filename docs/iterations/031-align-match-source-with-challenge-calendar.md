# Iteration 031 — Align Match Source With Challenge Calendar

## Context

Runtime validation on Raspberry Pi revealed that the engine is successfully connecting to MPP and generating opportunities.

However, the detected opportunities are based on matches such as:

- AS Saint-Étienne
- Guingamp
- Troyes
- Red Star
- Nancy

while the active challenge is:

```text
mpp_challenge_UDKDDH27
Mondial coin coin
```

which is attached to:

```json
{
  "championshipId": 8,
  "code": "cdm",
  "brand": "CDM 2026",
  "name": "World Cup 26"
}
```

The challenge UI and MPP API both confirm that the challenge is currently running on the 2026 World Cup competition.

Current Game Week:

```json
{
  "gameWeekNumber": 2,
  "startDate": "2026-06-18T16:00:00.000Z",
  "endDate": "2026-06-24T02:00:00.000Z"
}
```

The engine currently retrieves matches through:

```go
/championships-current-matches
```

and no usage of:

```text
/championship-calendar
/championship-match/summaries
```

was found in the codebase.

This creates a strong suspicion that opportunities are generated from a global match feed rather than the matches belonging to the challenge.

---

## Goal

Ensure that all opportunity detection is based exclusively on the matches belonging to the active challenge.

---

## Expected Match Retrieval Flow

Current challenge:

```text
challengeId
↓
GET /challenge/{challengeId}
↓
championshipId
↓
GET /championship-calendar/{championshipId}/nearest-game-weeks
↓
currentGameWeek.matchesIds
↓
POST /championship-match/summaries
↓
MatchSummary[]
↓
Opportunity Engine
```

---

## Investigation Tasks

### 1. Audit current match retrieval

Identify:

- where match IDs originate
- where MatchSummary objects are built
- where championships-current-matches is consumed
- whether challenge scope is lost during retrieval

Produce:

```text
docs/match-source-audit.md
```

---

### 2. Verify challenge consistency

Confirm:

```text
challenge
→ championshipId
→ currentGameWeek
→ matchesIds
```

is sufficient to reconstruct all matches required by the engine.

---

### 3. Implement challenge-scoped retrieval

Replace global match retrieval by:

```text
challenge
→ championshipId
→ nearest-game-weeks
→ currentGameWeek
→ match summaries
```

if investigation confirms current implementation is incorrect.

---

### 4. Validate opportunity generation

Run dry-run validation and confirm that generated opportunities reference World Cup teams such as:

- Canada
- Qatar
- Mexico
- South Korea
- USA
- Australia
- Czech Republic
- South Africa

instead of unrelated club competitions.

---

### 5. Add regression tests

Add tests guaranteeing that:

- challenge match retrieval uses challenge championshipId
- current game week is used
- match summaries originate from challenge matches
- opportunities cannot be generated from unrelated competitions

---

## Deliverables

- docs/match-source-audit.md
- updated MPP client if required
- updated tests
- updated architecture documentation if match retrieval changes

---

## Definition of Done

- Match source fully documented
- Challenge scope preserved end-to-end
- World Cup challenge produces World Cup opportunities
- No opportunity generated from unrelated competitions
- go test ./...
- go vet ./...
- git diff --check
- Worktree clean