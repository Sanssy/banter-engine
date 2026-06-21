# Iteration 026 - Opportunity Catalog Runtime Integration

## Goal

Make the Opportunity Catalog the single source of truth for:

- Detection
- Banter generation
- LLM prompting
- Future UI/API exposure

---

## Motivation

The project already contains:

- Opportunity detectors
- Banter generation
- Opportunity documentation

The catalog must become executable data instead of documentation only.

---

## New Structure

resources/

├── opportunities.json
├── archetypes.json
└── narratives.json

---

## opportunities.json

Each opportunity definition contains:

{
  "id": "RankingOvertake",
  "category": "Ranking",
  "severity": 2,
  "description": "A player overtakes another player in standings",
  "tags": ["ranking", "leaderboard"]
}

---

## Domain Model

type OpportunityDefinition struct {
    ID string
    Category string
    Severity int
    Description string
    Tags []string
}

---

## Catalog Loader

Create:

internal/catalog

Expose:

LoadOpportunityCatalog(path string) ([]OpportunityDefinition, error)

---

## Opportunity Validation

When an opportunity is emitted:

Opportunity{
    Type: "RankingOvertake"
}

The engine must validate that the opportunity exists in opportunities.json.

Unknown opportunities must generate an error.

---

## LLM Context Builder

Create:

internal/context

Expose:

BuildLLMContext(
    opportunity Opportunity,
    definition OpportunityDefinition,
) string

The generated prompt must include:

- opportunity id
- category
- description
- severity
- metadata

---

## Banter Generator

Replace hardcoded opportunity descriptions.

Use catalog definitions instead.

The generator must support:

- templates
- future LLM providers

---

## Future Compatibility

The catalog must support:

- 100+ opportunities
- player archetypes
- season narratives

without code modifications.

---

## Deliverables

Implement:

- resources/opportunities.json
- internal/catalog
- internal/context

Integrate catalog loading into startup.

Integrate catalog lookup into banter generation.

Integrate catalog information into LLM prompt generation.

---

## Definition Of Done

The application starts successfully and loads:

resources/opportunities.json

An emitted opportunity automatically gains access to:

- category
- severity
- description
- tags

without hardcoded mappings.
