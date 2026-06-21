# Club Name Resolution Diagnostic

## Run context

- Challenge: `mpp_challenge_UDKDDH27`
- Championship: `8`
- Current game week: `2`
- Matches: `24`
- Club references decoded from `GET /championship-clubs`: `138`

## Findings

The club payload is decoded correctly for references it contains. The first club in sorted order was:

```text
id=mpp_championship_club_1
name=Man. United
languages=[en-GB es-ES fr-FR]
```

The club identifiers used by the World Cup match summaries are not present in the decoded `/championship-clubs` map. For example:

```text
match=mpp_championship_match_2608266
homeClubId=mpp_championship_club_367
home_present=false
awayClubId=mpp_championship_club_522
away_present=false
```

The same result was observed for the first five matches, including club IDs `497`, `537`, `597`, `1873`, `659`, `1041`, `596`, and `575`.

## Context comparison

Without an application context, `GET /championship-clubs` returns 138 references and omits the championship `8` clubs.

With the OpenAPI-documented header below, the endpoint returns 148 references and includes the required international clubs:

```http
app-context: internationalEvent
```

Examples from the contextualized response:

```text
mpp_championship_club_522 -> Afrique du Sud
mpp_championship_club_614 -> Brésil
```

The objects also contain their championship `8` assets, jerseys, group data, and localized names.

## Conclusion

The failure is not caused by an ID-format mismatch, resolver lookup logic, or localized-name decoding. The match summaries and resolver use the same complete IDs.

The missing `app-context: internationalEvent` request header causes `/championship-clubs` to return the wrong reference set. A later corrective ticket should add this context when loading international club references. No functional correction was applied in this diagnostic iteration.
