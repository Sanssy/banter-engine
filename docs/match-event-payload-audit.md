# MPP Match Event Payload Audit

## Failure

The match detail endpoint returned an `eventsTimeline.score` JSON string, while `matchEventDTO` required an object:

```json
{"home": 1, "away": 0}
```

Go therefore stopped the run with `cannot unmarshal string into Go struct field matchEventDTO.eventsTimeline.score`. The local OpenAPI document defines `eventsTimeline` only as an array of untyped objects and does not constrain this field.

## Compatibility

`eventScoreDTO` now accepts both the historical object and string forms such as:

```json
"1-0"
"1 - 0"
"1:0"
```

An unrecognized string no longer aborts event retrieval; its raw value is logged and the domain score remains zero. Event type, time, side, player, and identifier continue to be decoded normally.

## Temporary Diagnostics

For each match event request, the client logs:

- the exact request URL;
- the first 2048 bytes of the raw response;
- the raw `score` value for the first five timeline events.

These diagnostics are intentionally bounded but may still contain match and player identifiers. They should be removed after the Raspberry dry-run confirms the live payload variants.
