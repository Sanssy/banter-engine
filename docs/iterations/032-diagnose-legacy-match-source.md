# Ticket 032 - Identifier pourquoi le moteur analyse encore les matchs 2569xxx

## Contexte

Le ticket 031 a été implémenté. Les vérifications API confirment que le challenge `mpp_challenge_UDKDDH27`, le championnat 8, la game week 2 et ses summaries utilisent tous des IDs `2608xxx`.

Le moteur continue pourtant à produire des matchs de clubs associés aux IDs `2569xxx`. Supprimer intégralement le dossier `data` ne change pas le résultat : les snapshots sont recréés avec ces anciens matchs.

## Travail demandé

Identifier précisément la source de données qui injecte les matchs `2569xxx`.

Ajouter des logs temporaires :

- après récupération des IDs de la game week ;
- après récupération des summaries ;
- avant l'exécution des détecteurs ;
- avant l'écriture de `data/matches.json` ;
- pour les routes MPP utilisées et l'identité du binaire exécuté.

Rechercher les fallbacks historiques, mocks, données de bootstrap, anciens providers, anciens chemins d'exécution et toute reconstruction depuis les forecasts.

## Résultat attendu

Produire un rapport expliquant :

1. où les matchs `2569xxx` sont injectés ;
2. pourquoi ils survivent au ticket 031 ;
3. quelle correction appliquer ultérieurement.

Aucune correction fonctionnelle dans ce ticket. Le diagnostic doit être isolé dans un commit séparé.
