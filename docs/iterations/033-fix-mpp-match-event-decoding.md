# Ticket 033 - Corriger le décodage des événements de match MPP

## Contexte

Le moteur récupère correctement les matchs de la Coupe du Monde, mais le run échoue lors du décodage des événements :

```text
json: cannot unmarshal string into Go struct field matchEventDTO.eventsTimeline.score
```

## Travail demandé

- journaliser temporairement la route, un extrait du JSON brut et les champs `eventsTimeline.score` ;
- comparer le payload avec le DTO actuel ;
- supporter le format réellement renvoyé par MPP ;
- conserver les tests existants ;
- ajouter un test reproduisant un score encodé sous forme de chaîne.

## Validation

Le dry-run doit charger les matchs `260826x`, récupérer leurs événements et terminer sans erreur de désérialisation.
