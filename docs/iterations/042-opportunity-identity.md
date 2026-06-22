# Ticket 042 - Opportunity Identity

## Contexte
Les opportunités ne possèdent pas d'identité métier stable.

## Objectif
Préparer le moteur à la déduplication avancée, aux statistiques et aux futurs digests.

## Évolution

```go
type Opportunity struct {
    Type    OpportunityType
    Actor   string
    Target  string
    MatchID string
    EventID string
    Key     string
}
```

## Critères d'acceptation
- Chaque opportunité possède une clé stable.
- Les opportunités liées à un match embarquent MatchID.
- Aucun comportement existant n'est cassé.
