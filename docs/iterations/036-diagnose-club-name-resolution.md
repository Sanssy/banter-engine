# Iteration 036 - Diagnose Club Name Resolution

## Goal

Identify why club identifiers loaded from MPP are not resolved and are displayed as `Equipe inconnue`.

## Diagnostic instrumentation

- Log the first club reference in deterministic ID order, including its resolved display name and available languages.
- Log the home and away club identifiers for the first five match summaries.
- Log whether those club identifiers exist in the decoded DTO, with their calculated names and languages.
- Log every club lookup with the requested ID, lookup status, and returned name.
- On lookup failure, log up to ten sorted keys registered in the resolver.

## Scope

This iteration adds diagnostic instrumentation only. It does not change club-name resolution behavior.
