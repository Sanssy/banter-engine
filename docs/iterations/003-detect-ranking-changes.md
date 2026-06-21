# Iteration 003 - Detect Ranking Changes

## Goal
Compare previous and current standings.

## Requirements
- Create internal/opportunities package
- Create Opportunity model:
```go
type Opportunity struct {
    Type string
    Actor string
    Target string
}
```
- Implement RankingOvertake detection

## Definition of Done
Detect when one user overtakes another in ranking.
