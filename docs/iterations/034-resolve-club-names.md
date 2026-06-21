# Ticket 034 - Résolution des noms de clubs

## Contexte

Le moteur utilise les bons matchs de la Coupe du Monde mais affiche encore des identifiants techniques `mpp_championship_club_*` dans les opportunités.

## Travail demandé

- vérifier le mapping de `/championship-clubs` ;
- tracer `MatchSummary -> clubId -> clubName -> Detector -> Opportunity` ;
- identifier la perte du nom ;
- garantir qu'aucun message utilisateur ne contient le préfixe technique.

## Hors périmètre

La récupération des matchs, le calcul des opportunités, les snapshots et le parsing des événements ne changent pas.
