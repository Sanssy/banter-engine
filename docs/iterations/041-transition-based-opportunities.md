# Ticket 041 - Transition-Based Opportunities

## Contexte
Certaines opportunités sont recalculées à chaque cycle tant que leur condition reste vraie.
Elles sont donc réémises plusieurs fois pour le même match.

## Objectif
Déclencher ces opportunités uniquement lors d'une transition métier significative.

## À corriger
- HugeUpset
- EveryoneWasWrong
- PredictionMassacre
- PopularMistake
- TheChosenOne
- AgainstTheCrowd
- CrowdFavorite
- CrowdTrap
- HotStreak
- ColdStreak
- LastPlaceLocked

## Règle
Une opportunité doit être générée lors du passage false -> true
ou lors d'un franchissement de seuil.

## Critères d'acceptation
- Plus aucune réémission d'une opportunité de résultat final sur les runs suivants.
- Les opportunités live continuent de fonctionner.
- Les tests couvrent les transitions.
