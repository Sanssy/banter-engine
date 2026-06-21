# Iteration 029 - Cleanup & Architecture Hardening

## Goal

Prepare the project for Raspberry deployment and future Ollama integration.

Before performing any cleanup, resolve all inconsistencies between the runtime opportunity catalog and the active detectors.

---

## Phase 1 - Opportunity Reconciliation (MANDATORY)

### Problem

Several detectors emit opportunities that are absent from the runtime catalog.

This can lead to:

- unknown opportunity errors
- interrupted execution cycles
- inconsistent runtime behavior

### Task

Produce a report containing:

```text
Detector
→ Opportunity emitted
→ Present in catalog? (yes/no)
```

Example:

```text
ranking.go
→ RankingOvertake
→ yes

streaks.go
→ HotStreak
→ no

surprises.go
→ HugeUpset
→ no
```

### Expected Output

Create:

`docs/opportunity-reconciliation-report.md`

containing:

- all emitted opportunities
- catalog presence status
- detector source file

### Decision Rule

DO NOT remove opportunities from detectors.

Instead:

- treat detectors as the source of truth
- extend the catalog when opportunities are missing

If detectors emit 113 opportunities:

```text
catalog size = 113
```

The catalog must represent reality.

The previous target of exactly 100 opportunities is no longer a constraint.

---

## Phase 2 - Catalog Update

Add all missing opportunities to:

`resources/opportunities.json`

Every added opportunity must follow the official template:

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

Validate:

- unique ids
- valid categories
- no missing required fields

---

## Phase 3 - LLM Layer Review

Inspect:

```text
internal/banter/llm.go
internal/banter/provider.go
internal/context/
internal/contextbuilder/
```

If unused:

- remove them

If referenced:

- document why they must remain

Produce a short justification in:

`docs/llm-review.md`

---

## Phase 4 - Catalog Naming Cleanup

Review:

```text
internal/catalog/
internal/opportunities/catalog.go
```

Responsibilities appear distinct.

Do not merge automatically.

If opportunities/catalog.go only contains domain identifiers or definitions:

- rename to a more explicit name

Example:

```text
opportunity_definitions.go
opportunity_registry.go
```

---

## Phase 5 - Repository Cleanup

Remove:

```text
.DS_Store
__MACOSX/
```

Add or update:

`.gitignore`

Verify:

- no generated files committed
- no platform-specific artifacts

---

## Phase 6 - Documentation Cleanup

Review and remove obsolete files:

```text
docs/opportunities-1000.json
docs/resources/opportunities.json
```

Only remove if confirmed unused.

Also review:

```text
resources/archetypes.json
resources/narratives.json
```

If unused:

- remove them
OR
- document future usage

---

## Phase 7 - Dependency Review

Verify:

- every package is used
- no dead code remains
- no unused exported APIs remain

---

## Definition of Done

- all detector opportunities reconciled
- no unknown opportunity possible
- runtime catalog updated
- dead code removed
- obsolete files removed
- naming clarified
- .gitignore cleaned

Validation:

```bash
go test ./...
go vet ./...
git diff --check
```

Produce a final summary:

```text
What was removed
What was renamed
What was added
Remaining technical debt
Ready for Raspberry deployment? (yes/no)
Ready for Ollama integration? (yes/no)
```
