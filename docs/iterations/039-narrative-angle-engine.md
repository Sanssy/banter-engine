# Ticket 039 - Narrative Angle Engine

## Contexte

Les premiers tests Ollama montrent que la qualité dépend davantage de l'angle narratif fourni au modèle que du modèle lui-même.

Exemple :

Opportunity:
- Paraguay bat la Turquie
- 92% avaient choisi la Turquie

Un LLM peut produire :
- "Quel choc !" (faible)
- "Les 8% restants paradent déjà." (fort)

La différence provient de l'angle narratif utilisé.

---

## Objectif

Introduire une couche intermédiaire :

Opportunity
→ Narrative Angle
→ Narrator
→ Message

---

## Narrative Angles initiaux

### CrowdWrong
La majorité s'est trompée.

Utilisé pour :
- EveryoneWasWrong
- PredictionMassacre
- PopularPickCollapse

### MinorityVictory
Une minorité avait raison.

Utilisé pour :
- TheChosenOne
- AgainstTheCrowd

### FallFromGrace
Chute d'un favori ou d'un joueur.

### Rise
Progression spectaculaire.

### Curse
Série noire.

### Dominance
Domination ou prise de pouvoir.

---

## Critères d'acceptation

- chaque opportunité possède un angle narratif
- le narrateur reçoit l'angle en plus des faits
- aucun changement fonctionnel visible côté utilisateur
- préparation du terrain pour les prompts spécialisés par angle
