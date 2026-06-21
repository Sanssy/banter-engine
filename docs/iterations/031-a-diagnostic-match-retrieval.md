# Ticket 031-A - Diagnostic récupération des matchs

## Objectif

Comprendre pourquoi le moteur continue à analyser les matchs historiques `2569xxx` alors que le challenge est rattaché au championnat Coupe du Monde 2026 (`championshipId=8`).

## Constat

- `challengeId = mpp_challenge_UDKDDH27`
- `championshipId = 8`
- `currentGameWeek = 2`
- `nearest-game-weeks` renvoie des matchs `2608265` à `2608288`
- le moteur analyse pourtant des matchs `2569036`, `2569037`, etc.
- vider complètement le dossier `data` ne change rien

## Travail demandé

Ajouter des logs temporaires détaillés dans le pipeline de récupération des matchs.

Afficher :

1. `challengeId` utilisé
2. `championshipId` utilisé
3. game week sélectionnée
4. nombre de `matchIds` récupérés
5. les 10 premiers `matchIds` récupérés
6. les 5 premiers matchs après résolution des summaries :
   - `MatchID`
   - `HomeTeam`
   - `AwayTeam`
   - `Date`
7. la route MPP appelée à chaque étape

## Objectif final

Identifier précisément à quel moment les matchs `2608xxx` deviennent des matchs `2569xxx`.

Aucune correction fonctionnelle. Uniquement instrumentation et diagnostic.
