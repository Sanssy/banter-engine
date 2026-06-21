package banter

import (
	"context"
	"fmt"
	"strings"

	"github.com/DSanoussy/banter-engine/internal/opportunities"
)

type Provider interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type Generator struct {
	provider Provider
}

func NewGenerator(provider Provider) *Generator {
	return &Generator{provider: provider}
}

func (g *Generator) Generate(ctx context.Context, op opportunities.Opportunity) string {
	if g.provider != nil {
		message, err := g.provider.Generate(ctx, opportunityPrompt(op))
		if err == nil && strings.TrimSpace(message) != "" {
			return strings.TrimSpace(message)
		}
	}

	return Generate(op)
}

func opportunityPrompt(op opportunities.Opportunity) string {
	return fmt.Sprintf(
		"Rédige un message de banter footballistique court en français. Type: %s. Acteur: %s. Cible: %s.",
		op.Type,
		op.Actor,
		op.Target,
	)
}
