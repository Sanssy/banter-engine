# Ticket 040 - Curated Few-Shot Library

## Contexte

Gemma4:E2B réagit beaucoup mieux lorsqu'on lui fournit plusieurs exemples cohérents du style Banter Engine.

L'objectif n'est pas d'envoyer des centaines d'exemples mais de sélectionner les plus pertinents.

---

## Objectif

Créer une bibliothèque d'exemples narratifs versionnés.

Structure :

resources/narratives/examples.json

Chaque exemple contient :

- category
- angle
- facts
- message

---

## Runtime

Avant un appel LLM :

Opportunity
→ recherche des meilleurs exemples
→ injection de 3 à 5 exemples
→ génération

---

## Critères d'acceptation

- bibliothèque chargée au démarrage
- sélection des exemples compatibles avec l'opportunité
- maximum 5 exemples injectés
- configuration simple sans base de données
