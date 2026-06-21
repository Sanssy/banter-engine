# Iteration 038 - Résumé intelligent et stratégie de diffusion des notifications

## Objectif

Réduire le bruit des notifications en filtrant les opportunités par importance, en instaurant une plage silencieuse nocturne et en générant un récapitulatif matinal.

---

## Ce qui a changé

### Nouveau package `internal/notify`

Trois fichiers créés.

#### `schedule.go`

```
IsQuietHour(t)         → true entre 00h00 et 07h59
IsNightSummaryHour(t)  → true à 08h00 exactement
```

#### `selector.go`

```
SelectTop(ops, catalog, n)
```

- trie les opportunités par `catalog.Severity` décroissant
- déduplique par type (une seule opportunité par type)
- conserve au maximum `n` résultats
- constante `MaxNotificationsPerRun = 5`

#### `summarizer.go`

```
NightSummary(ops, catalog)  → texte multi-lignes pour le digest matinal
LiveDigest(ops, catalog)    → non utilisé directement (utilitaire)
```

Format du digest :

```
**Pendant la nuit :**
🔥 Sandrine — BiggestWinner
📉 MomoDuParisFc — FreeFall
📈 Luc_arnes vs ChaOli — RankingOvertake
```

---

### Snapshot

Deux nouvelles paires de fonctions dans `internal/snapshot/snapshot.go` :

- `SaveNightBuffer` / `LoadNightBuffer` → `night_buffer.json`
  Accumule les opportunités détectées pendant la plage silencieuse.

- `SaveNightSummaryDate` / `LoadNightSummaryDate` → `night_summary_date.txt`
  Stocke la date du dernier digest envoyé pour éviter les doublons.

---

### Engine

`runOnce()` délègue maintenant à trois chemins distincts selon l'heure :

| Heure | Comportement |
|---|---|
| 00h00 – 07h59 | Accumulation dans `night_buffer.json`, zéro notification |
| 08h00 (premier run du jour) | Envoi du digest de nuit + dispatch normal filtré |
| 08h01 – 23h59 | Dispatch normal filtré (top 5 par sévérité) |

La résolution des noms Actor/Target est désormais effectuée avant le filtrage, pour que le digest contienne les vrais noms.

---

## Fichiers modifiés

```
internal/notify/schedule.go       (nouveau)
internal/notify/selector.go       (nouveau)
internal/notify/summarizer.go     (nouveau)
internal/snapshot/snapshot.go     (étendu)
internal/engine/engine.go         (mis à jour)
```

---

## Critères d'acceptation

- [x] aucune notification entre 00h00 et 07h59
- [x] récapitulatif de nuit envoyé à 08h00 (une fois par jour)
- [x] opportunités filtrées avant diffusion (top 5 par sévérité)
- [x] le moteur conserve toutes les opportunités détectées (buffer persisté)
- [x] compatible avec une future synthèse LLM (summarizer remplaçable)

---

## Validation

```bash
go build ./...
go test ./...
```

109 tests passent, 0 régression.
