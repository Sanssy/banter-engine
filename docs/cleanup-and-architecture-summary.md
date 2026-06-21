# Cleanup and Architecture Summary

## What Was Removed

- The unused LLM provider abstraction and its tests.
- The unused `internal/context` and `internal/contextbuilder` packages.
- Placeholder `resources/archetypes.json` and `resources/narratives.json` files with no runtime consumer.
- Obsolete generated catalogs under `docs/opportunities-1000.json` and `docs/resources/opportunities.json`.
- Superseded placeholder documentation files `docs/04-opportunity-catalog.md` through `docs/08-season-narratives.md`.
- macOS `.DS_Store` files.
- Unnecessarily exported detector helpers and the exported point-impact implementation type.

## What Was Renamed

- `internal/opportunities/catalog.go` became `internal/opportunities/opportunity_registry.go` to reflect that it owns domain opportunity identifiers, not catalog loading.

## What Was Added

- Thirteen missing detector opportunity definitions, bringing the runtime catalog to 113 entries.
- An explicit registry of the 37 opportunity types emitted by active detectors.
- Startup and test validation that every registered detector opportunity exists in the runtime catalog.
- The opportunity reconciliation and LLM review reports.
- A `.gitignore` covering macOS metadata, local runtime state, Go test binaries, profiles, and the local executable.

## Remaining Technical Debt

- MPP and Discord HTTP methods do not accept a `context.Context`; graceful shutdown can therefore wait for their configured HTTP timeouts.
- The runtime catalog path is relative to the working directory. The documented systemd unit sets the required working directory.
- Ollama configuration, transport, and generation behavior remain intentionally unimplemented.

## Readiness

- Ready for Raspberry deployment? **Yes**, subject to the documented environment configuration and an on-device smoke test.
- Ready for Ollama integration? **Yes**, as a clean starting point; no Ollama runtime is included in this iteration.
