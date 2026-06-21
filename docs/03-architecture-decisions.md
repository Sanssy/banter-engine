# Architecture Decisions

## Challenge-Scoped Match Source

**Decision:** Match IDs used by the opportunity engine must originate from the configured challenge calendar.

The runtime resolves matches through:

```text
challenge ID
-> challenge championship ID
-> nearest game weeks
-> current game week match IDs
-> batch match summaries
```

The global `/championships-current-matches` feed must not be used for opportunity detection because it combines unrelated competitions. Responses from the summaries endpoint are filtered against the requested calendar IDs before conversion to domain matches.

Challenge scope is preserved from configuration through match retrieval, forecasts, snapshots, and opportunity generation. The global club directory may be used for display-name enrichment only; it cannot provide match identity.
