# Iteration 002 - Save Snapshot

## Goal
Persist standings locally.

## Requirements
- Create internal/snapshot package
- Implement:
  - SaveStandings(path string, standings []model.Standing) error
  - LoadStandings(path string) ([]model.Standing, error)
- Store data in JSON
- If file does not exist, return empty standings and no error

## Structure
data/
└── standings.json

## Definition of Done
Current standings can be saved and loaded successfully.
