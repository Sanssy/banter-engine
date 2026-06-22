package narrator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sanssy/banter-engine/internal/catalog"
	"github.com/Sanssy/banter-engine/internal/narrative"
	"github.com/Sanssy/banter-engine/internal/notify"
	"github.com/Sanssy/banter-engine/internal/opportunities"
)

// OllamaNarrator calls a local Ollama instance and falls back to DeterministicNarrator on any failure.
type OllamaNarrator struct {
	url      string
	model    string
	timeout  time.Duration
	client   *http.Client
	examples *narrative.Library
	fallback DeterministicNarrator
}

func NewOllamaNarrator(url, model string, timeout time.Duration, examples *narrative.Library) *OllamaNarrator {
	return &OllamaNarrator{
		url:      strings.TrimRight(url, "/"),
		model:    model,
		timeout:  timeout,
		client:   &http.Client{Timeout: timeout},
		examples: examples,
	}
}

func (o *OllamaNarrator) Narrate(op opportunities.Opportunity, def catalog.OpportunityDefinition, angle narrative.Angle) string {
	examples := o.examples.Select(def.Category, angle, narrative.MaxExamples)
	prompt := buildLivePrompt(op, def, angle, examples)
	result, err := o.generate(prompt)
	if err != nil || strings.TrimSpace(result) == "" {
		return o.fallback.Narrate(op, def, angle)
	}
	return result
}

func (o *OllamaNarrator) Summarize(ops []opportunities.Opportunity, cat *catalog.Catalog) string {
	prompt := buildDigestPrompt(ops, cat, o.selectDigestExamples(ops, cat))
	result, err := o.generate(prompt)
	if err != nil || strings.TrimSpace(result) == "" {
		return o.fallback.Summarize(ops, cat)
	}
	return result
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func (o *OllamaNarrator) generate(prompt string) (string, error) {
	body, err := json.Marshal(ollamaRequest{
		Model:  o.model,
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return "", err
	}

	resp, err := o.client.Post(o.url+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return strings.TrimSpace(result.Response), nil
}

func buildLivePrompt(
	op opportunities.Opportunity,
	def catalog.OpportunityDefinition,
	angle narrative.Angle,
	examples []narrative.Example,
) string {
	var sb strings.Builder
	sb.WriteString("Tu es un commentateur de ligue de pronostics sportifs. ")
	sb.WriteString("Reformule l'événement suivant en une phrase courte, naturelle et légèrement amusante en français. ")
	sb.WriteString("Maximum 150 caractères. Une seule phrase. N'invente aucun fait.\n\n")
	sb.WriteString(fmt.Sprintf("Événement : %s\n", def.Name))
	if def.Description != "" {
		sb.WriteString(fmt.Sprintf("Description : %s\n", def.Description))
	}
	if op.Actor != "" {
		sb.WriteString(fmt.Sprintf("Acteur principal : %s\n", op.Actor))
	}
	if op.Target != "" {
		sb.WriteString(fmt.Sprintf("Cible / contexte : %s\n", op.Target))
	}
	sb.WriteString(fmt.Sprintf("Angle narratif : %s\n", angle))
	if guidance := angle.Guidance(); guidance != "" {
		sb.WriteString(fmt.Sprintf("Intention narrative : %s\n", guidance))
	}
	writeFewShotExamples(&sb, examples)
	sb.WriteString("\nRéponds uniquement avec la phrase reformulée, sans explication.")
	return sb.String()
}

func buildDigestPrompt(
	ops []opportunities.Opportunity,
	cat *catalog.Catalog,
	examples []narrative.Example,
) string {
	selected := notify.SelectTop(ops, cat, 8)

	var sb strings.Builder
	sb.WriteString("Tu es un commentateur de ligue de pronostics sportifs. ")
	sb.WriteString("Rédige un digest du matin en français pour raconter ce qu'il s'est passé cette nuit. ")
	sb.WriteString("Utilise des emojis pertinents. Sois informatif, bref et légèrement divertissant. ")
	sb.WriteString("N'invente aucun fait, score ou classement.\n\n")
	sb.WriteString("Événements de la nuit :\n")

	for _, op := range selected {
		def, found := cat.FindByID(op.Type)
		if !found {
			continue
		}
		line := fmt.Sprintf("- %s", def.Name)
		if op.Actor != "" {
			line += fmt.Sprintf(" : %s", op.Actor)
		}
		if op.Target != "" {
			line += fmt.Sprintf(" / %s", op.Target)
		}
		sb.WriteString(line + "\n")
	}

	writeFewShotExamples(&sb, examples)
	sb.WriteString("\nRéponds uniquement avec le digest mis en forme, sans explication.")
	return sb.String()
}

func (o *OllamaNarrator) selectDigestExamples(
	ops []opportunities.Opportunity,
	cat *catalog.Catalog,
) []narrative.Example {
	selectedOps := notify.SelectTop(ops, cat, 8)
	pools := make([][]narrative.Example, 0, len(selectedOps))
	for _, op := range selectedOps {
		definition, found := cat.FindByID(op.Type)
		if !found {
			continue
		}
		angle := narrative.ForOpportunity(op.Type)
		candidates := o.examples.Select(definition.Category, angle, narrative.MaxExamples)
		if len(candidates) > 0 {
			pools = append(pools, candidates)
		}
	}

	result := make([]narrative.Example, 0, narrative.MaxExamples)
	seen := make(map[string]struct{}, narrative.MaxExamples)
	for candidateIndex := 0; candidateIndex < narrative.MaxExamples; candidateIndex++ {
		for _, candidates := range pools {
			if candidateIndex >= len(candidates) {
				continue
			}
			example := candidates[candidateIndex]
			key := string(example.Angle) + "\x00" + example.Facts
			if _, duplicate := seen[key]; duplicate {
				continue
			}
			seen[key] = struct{}{}
			result = append(result, example)
			if len(result) == narrative.MaxExamples {
				return result
			}
		}
	}
	return result
}

func writeFewShotExamples(sb *strings.Builder, examples []narrative.Example) {
	if len(examples) == 0 {
		return
	}
	if len(examples) > narrative.MaxExamples {
		examples = examples[:narrative.MaxExamples]
	}

	sb.WriteString("\nExemples de style à suivre sans les copier mot pour mot :\n")
	for index, example := range examples {
		sb.WriteString(fmt.Sprintf(
			"Exemple %d - faits : %s\nExemple %d - message : %s\n",
			index+1,
			example.Facts,
			index+1,
			example.Message,
		))
	}
}
