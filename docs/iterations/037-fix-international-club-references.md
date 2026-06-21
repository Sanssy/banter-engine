# Iteration 037 - Fix International Club References

## Goal

Resolve championship `8` club identifiers to their localized names.

## Root cause

`GET /championship-clubs` returns the standard club reference set when no application context is supplied. World Cup clubs are returned when the request includes:

```http
app-context: internationalEvent
```

## Implementation

- Send the international event context only when loading `/championship-clubs`.
- Keep all other MPP requests unchanged.
- Verify the header and localized club-name mapping in the MPP client test.
