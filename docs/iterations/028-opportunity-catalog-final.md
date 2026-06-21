# Iteration 028 - Opportunity Catalog Final

## Goal

Build the definitive opportunity catalog used by Banter Engine.

## Context

The project now contains 100 documented opportunities distributed across:

- Ranking
- Predictions
- Crowd
- MatchEvents
- Narratives

The objective is to consolidate them into a single runtime catalog.

---

## Scope

Create:

resources/opportunities.json

The catalog becomes the single source of truth used by detectors,
context builders and future LLM integrations.

---

## Opportunity Schema

Each opportunity MUST follow:

```yaml
id:
name:
category:
severity:
description:
requiredData:
trigger:
banterAngles:
relatedOpportunities:
tags:
```

---

## Tasks

### 1. Consolidate Catalog

Merge all validated packs into:

resources/opportunities.json

Expected quantity:

100 opportunities

---

### 2. Validate Catalog

Implement validation rules:

- unique ids
- valid category
- severity between 1 and 5
- non-empty description
- non-empty tags

Startup must fail if catalog is invalid.

---

### 3. Catalog Loader

Create:

internal/catalog/

Expose:

LoadCatalog(path string)

Responsibilities:

- load JSON
- validate schema
- expose lookup APIs

---

### 4. Catalog Queries

Support:

FindByID(id)

FindByCategory(category)

FindRelated(id)

---

### 5. Tests

Add tests covering:

- valid catalog
- invalid catalog
- duplicate ids
- missing fields

---

## Non Goals

Do NOT introduce:

- Ollama
- LLM prompts
- Discord changes
- new opportunity types

---

## Definition Of Done

- opportunities.json exists
- 100 opportunities loaded
- validation implemented
- tests passing
- go test ./...
- go vet ./...
