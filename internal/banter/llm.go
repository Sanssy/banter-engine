package banter

import (
	"context"
	"fmt"
	"strings"

	"github.com/DSanoussy/banter-engine/internal/catalog"
	runtimecontext "github.com/DSanoussy/banter-engine/internal/context"
	"github.com/DSanoussy/banter-engine/internal/contextbuilder"
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
	return g.GenerateWithContext(ctx, op, contextbuilder.Context{})
}

func (g *Generator) GenerateWithContext(
	ctx context.Context,
	op opportunities.Opportunity,
	banterContext contextbuilder.Context,
) string {
	return g.GenerateWithDefinition(ctx, op, catalog.OpportunityDefinition{ID: op.Type}, banterContext)
}

func (g *Generator) GenerateWithDefinition(
	ctx context.Context,
	op opportunities.Opportunity,
	definition catalog.OpportunityDefinition,
	banterContext contextbuilder.Context,
) string {
	if g.provider != nil {
		message, err := g.provider.Generate(ctx, opportunityPrompt(op, definition, banterContext))
		if err == nil && strings.TrimSpace(message) != "" {
			return strings.TrimSpace(message)
		}
	}

	return GenerateWithDefinition(op, definition)
}

func opportunityPrompt(
	op opportunities.Opportunity,
	definition catalog.OpportunityDefinition,
	banterContext contextbuilder.Context,
) string {
	summary, err := banterContext.Summary()
	if err != nil {
		summary = "{}"
	}
	return fmt.Sprintf(
		"Rédige un message de banter footballistique court en français. Opportunité: %s. Contexte runtime: %s.",
		runtimecontext.BuildLLMContext(op, definition),
		summary,
	)
}
