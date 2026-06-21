# Iteration 037 - Cleanup MPP Investigation Logs

## Goal

Remove temporary MPP investigation logs while preserving concise operational visibility.

## Kept logs

- Build identity and scheduler lifecycle.
- Challenge, championship, game week, and match count.
- Aggregate club and user reference counts.
- Detector input count, snapshot persistence, and completed opportunity count.
- All existing errors.

## Removed logs

- Raw JSON previews.
- Per-club and per-user lookups.
- Per-match resolution details.
- Match event payloads and raw score traces.
- Request routes and summary previews added during diagnostics.
