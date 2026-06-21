// bench-narrator benchmarks Ollama models for Banter Engine narrative quality.
// Run on the target device: go run ./cmd/bench-narrator
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const ollamaURL = "http://localhost:11434"

var models = []string{
	"gemma4:e4b",
	"qwen3:4b",
	"phi4-mini",
}

type scenario struct {
	name   string
	prompt string
	maxLen int
	// weight controls how much this scenario influences the final recommendation.
	// The digest is the primary use case; live notifications are secondary.
	weight float64
	// runs is how many times to run this scenario (>1 reveals variety and consistency).
	runs int
}

var scenarios = []scenario{
	{
		name: "live_upset",
		prompt: `Tu es un commentateur de ligue de pronostics sportifs.
Reformule l'événement suivant en une phrase courte, naturelle et légèrement amusante en français.
Maximum 150 caractères. Une seule phrase. N'invente aucun fait.

Événement : Gros favori renversé
Acteur principal : Paraguay
Cible / contexte : Turquie

Réponds uniquement avec la phrase reformulée, sans explication.`,
		maxLen: 150,
		weight: 1.0,
		runs:   1,
	},
	{
		name: "live_overtake",
		prompt: `Tu es un commentateur de ligue de pronostics sportifs.
Reformule l'événement suivant en une phrase courte, naturelle et légèrement amusante en français.
Maximum 150 caractères. Une seule phrase. N'invente aucun fait. Pas de moquerie agressive.

Événement : Dépassement au classement
Acteur principal : Luc_arnes
Cible / contexte : ChaOli

Réponds uniquement avec la phrase reformulée, sans explication.`,
		maxLen: 150,
		weight: 1.0,
		runs:   1,
	},
	// The morning digest is the primary use case.
	// It is run 3 times to evaluate consistency and variety.
	// A user who did not watch the matches must understand what happened overnight,
	// who won, who lost, who moved up the standings — and close the message smiling.
	{
		name: "digest_morning",
		prompt: `Tu es un commentateur de ligue de pronostics sportifs.
Un utilisateur qui n'a pas regardé les matchs doit comprendre ce qui s'est passé pendant la nuit,
qui a gagné, qui a perdu, qui a monté au classement, et repartir avec le sourire.
Rédige un digest du matin en français avec des emojis pertinents.
Sois informatif, vivant et légèrement divertissant. N'invente aucun fait, score ou classement.
Maximum 800 caractères.

Événements de la nuit :
- Victoire écrasante : Espagne / Arabie Saoudite (4-0)
- Gros favori renversé : Turquie
- Gros favori renversé : Équateur
- Gros gain de points : Sandrine / +50 points
- Dépassement au classement : Luc_arnes / ChaOli
- Série de 5 échecs : LeDaveCoinCoin

Réponds uniquement avec le digest mis en forme, sans explication.`,
		maxLen: 800,
		weight: 3.0,
		runs:   3,
	},
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response           string `json:"response"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int64  `json:"prompt_eval_duration"`
	TotalDuration      int64  `json:"total_duration"`
}

type result struct {
	model    string
	scenario string
	text     string
	duration time.Duration
	tokens   int
	tokensPS float64
	ramMB    int
	lenOK    bool
	err      error
}

func main() {
	client := &http.Client{Timeout: 120 * time.Second}

	available := checkAvailableModels(client)
	if len(available) == 0 {
		fmt.Fprintln(os.Stderr, "Aucun modèle disponible — vérifiez qu'Ollama est démarré et que des modèles sont installés.")
		os.Exit(1)
	}

	fmt.Printf("Modèles disponibles : %s\n\n", strings.Join(available, ", "))

	var results []result

	for _, model := range available {
		fmt.Printf("=== %s ===\n", model)
		// Warm-up: one silent generation to load the model into RAM.
		warmup(client, model)
		ramMB := measureRAM(model)
		for _, sc := range scenarios {
			for run := 1; run <= sc.runs; run++ {
				r := runScenario(client, model, sc, ramMB)
				if sc.runs > 1 {
					r.scenario = fmt.Sprintf("%s #%d", sc.name, run)
				}
				results = append(results, r)
				printResult(r)
			}
		}
		fmt.Println()
	}

	printSummaryTable(results)
	printRecommendation(results)
}

func checkAvailableModels(client *http.Client) []string {
	resp, err := client.Get(ollamaURL + "/api/tags")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var payload struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil
	}

	installed := make(map[string]bool)
	for _, m := range payload.Models {
		installed[m.Name] = true
		// also index by short name (without tag)
		if short, _, found := strings.Cut(m.Name, ":"); found {
			installed[short] = true
		}
	}

	var available []string
	for _, m := range models {
		short, _, _ := strings.Cut(m, ":")
		if installed[m] || installed[short] {
			available = append(available, m)
		}
	}
	// Add any installed model not in our list that looks relevant.
	for _, m := range payload.Models {
		found := false
		for _, a := range available {
			if a == m.Name {
				found = true
				break
			}
		}
		if !found {
			available = append(available, m.Name)
		}
	}
	return available
}

func warmup(client *http.Client, model string) {
	fmt.Printf("  chargement %s...\n", model)
	body, _ := json.Marshal(ollamaRequest{Model: model, Prompt: "Bonjour", Stream: false})
	resp, err := client.Post(ollamaURL+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		return
	}
	resp.Body.Close()
}

func measureRAM(model string) int {
	// Use `ps` to find the ollama runner process and read its RSS.
	out, err := exec.Command("sh", "-c",
		"ps aux | grep -i 'ollama' | grep -v grep | awk '{sum += $6} END {print sum}'",
	).Output()
	if err != nil {
		return 0
	}
	kb, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	return kb / 1024
}

func runScenario(client *http.Client, model string, sc scenario, ramMB int) result {
	body, _ := json.Marshal(ollamaRequest{Model: model, Prompt: sc.prompt, Stream: false})

	start := time.Now()
	resp, err := client.Post(ollamaURL+"/api/generate", "application/json", bytes.NewReader(body))
	elapsed := time.Since(start)

	if err != nil {
		return result{model: model, scenario: sc.name, err: err, duration: elapsed, ramMB: ramMB}
	}
	defer resp.Body.Close()

	var olResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&olResp); err != nil {
		return result{model: model, scenario: sc.name, err: err, duration: elapsed, ramMB: ramMB}
	}

	text := strings.TrimSpace(olResp.Response)
	tokPS := 0.0
	if olResp.EvalDuration > 0 {
		tokPS = float64(olResp.EvalCount) / (float64(olResp.EvalDuration) / 1e9)
	}

	return result{
		model:    model,
		scenario: sc.name,
		text:     text,
		duration: elapsed,
		tokens:   olResp.EvalCount,
		tokensPS: tokPS,
		ramMB:    ramMB,
		lenOK:    len([]rune(text)) <= sc.maxLen,
	}
}

func printResult(r result) {
	fmt.Printf("\n  [%s]\n", r.scenario)
	if r.err != nil {
		fmt.Printf("  ERREUR: %v\n", r.err)
		return
	}
	fmt.Printf("  Durée    : %s\n", r.duration.Round(time.Millisecond))
	fmt.Printf("  Tokens/s : %.1f\n", r.tokensPS)
	fmt.Printf("  RAM      : ~%d Mo\n", r.ramMB)
	lenStatus := "✓"
	if !r.lenOK {
		lenStatus = "✗ TROP LONG"
	}
	fmt.Printf("  Longueur : %d chars %s\n", len([]rune(r.text)), lenStatus)
	fmt.Printf("  Réponse  :\n")
	// Indent the response
	for _, line := range strings.Split(r.text, "\n") {
		fmt.Printf("    %s\n", line)
	}
}

func printSummaryTable(results []result) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TABLEAU COMPARATIF")
	fmt.Println(strings.Repeat("=", 80))

	// Header
	fmt.Printf("%-20s %-22s %8s %8s %6s %6s\n",
		"Modèle", "Scénario", "Durée", "Tok/s", "RAM Mo", "Long.")
	fmt.Println(strings.Repeat("-", 72))

	for _, r := range results {
		status := "✓"
		if r.err != nil {
			status = "ERR"
		} else if !r.lenOK {
			status = "LONG"
		}
		dur := "-"
		tps := "-"
		ram := "-"
		if r.err == nil {
			dur = r.duration.Round(time.Millisecond).String()
			tps = fmt.Sprintf("%.1f", r.tokensPS)
			ram = fmt.Sprintf("%d", r.ramMB)
		}
		fmt.Printf("%-20s %-22s %8s %8s %6s %6s\n",
			r.model, r.scenario, dur, tps, ram, status)
	}
}

// scenarioWeight returns the weight defined in the scenarios list for a given scenario name.
// Digest runs are named "digest_morning #N" — they still match the digest weight.
func scenarioWeight(name string) float64 {
	for _, sc := range scenarios {
		if sc.name == name || strings.HasPrefix(name, sc.name) {
			return sc.weight
		}
	}
	return 1.0
}

func printRecommendation(results []result) {
	// Score each model: tokens/s × scenario_weight × (1 if lenOK else 0.5).
	// The digest has weight 3× so it dominates the recommendation — it is the primary use case.
	type modelScore struct {
		model  string
		score  float64
		avgTPS float64
		avgDur time.Duration
		errors int
	}

	scores := map[string]*modelScore{}
	counts := map[string]int{}

	for _, r := range results {
		if _, ok := scores[r.model]; !ok {
			scores[r.model] = &modelScore{model: r.model}
		}
		s := scores[r.model]
		counts[r.model]++
		if r.err != nil {
			s.errors++
			continue
		}
		lenWeight := 1.0
		if !r.lenOK {
			lenWeight = 0.5
		}
		w := scenarioWeight(r.scenario)
		s.score += r.tokensPS * w * lenWeight
		s.avgTPS += r.tokensPS
		s.avgDur += r.duration
	}

	var best *modelScore
	for _, s := range scores {
		n := counts[s.model]
		if n > 0 {
			s.avgTPS /= float64(n - s.errors)
			s.avgDur /= time.Duration(n)
		}
		if best == nil || s.score > best.score {
			best = s
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("RECOMMANDATION")
	fmt.Println(strings.Repeat("=", 80))

	if best == nil {
		fmt.Println("Impossible de déterminer un gagnant — tous les modèles ont échoué.")
		return
	}

	fmt.Printf("\nModèle retenu : %s\n", best.model)
	fmt.Printf("Vitesse moyenne : %.1f tokens/s\n", best.avgTPS)
	fmt.Printf("Durée moyenne   : %s\n", best.avgDur.Round(time.Millisecond))
	fmt.Println()
	fmt.Println("Configuration Ollama recommandée (.env) :")
	fmt.Printf("  OLLAMA_ENABLED=true\n")
	fmt.Printf("  OLLAMA_MODEL=%s\n", best.model)
	fmt.Printf("  OLLAMA_URL=http://localhost:11434\n")
	fmt.Printf("  OLLAMA_TIMEOUT=30s\n")
	fmt.Println()
	fmt.Println("Critère principal : qualité du digest du matin (poids ×3 dans le score).")
	fmt.Println("Les 3 runs du digest permettent d'évaluer la variété et la cohérence.")
	fmt.Println("Ajustez OLLAMA_TIMEOUT selon la vitesse observée sur le Pi.")
	fmt.Println("Sur Raspberry Pi 4 8 Go, comptez ~3-6 tokens/s pour un modèle 4B récent.")
}
