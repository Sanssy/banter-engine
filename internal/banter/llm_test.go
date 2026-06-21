package banter

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/DSanoussy/banter-engine/internal/opportunities"
)

func TestGeneratorUsesProvider(t *testing.T) {
	provider := &fakeProvider{message: "  Message généré par le provider.  "}
	generator := NewGenerator(provider)
	opportunity := opportunities.Opportunity{
		Type:   opportunities.RankingOvertake,
		Actor:  "Sanssy",
		Target: "William",
	}

	got := generator.Generate(context.Background(), opportunity)
	if got != "Message généré par le provider." {
		t.Fatalf("Generate() = %q, want provider message", got)
	}
	for _, value := range []string{opportunity.Type, opportunity.Actor, opportunity.Target} {
		if !strings.Contains(provider.prompt, value) {
			t.Fatalf("prompt %q does not contain %q", provider.prompt, value)
		}
	}
}

func TestGeneratorFallsBackToTemplate(t *testing.T) {
	opportunity := opportunities.Opportunity{
		Type:   opportunities.RankingOvertake,
		Actor:  "Sanssy",
		Target: "William",
	}
	want := Generate(opportunity)

	tests := []struct {
		name     string
		provider Provider
	}{
		{name: "no provider"},
		{name: "provider error", provider: &fakeProvider{err: errors.New("provider unavailable")}},
		{name: "empty response", provider: &fakeProvider{message: "  "}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := NewGenerator(tt.provider)
			if got := generator.Generate(context.Background(), opportunity); got != want {
				t.Fatalf("Generate() = %q, want template %q", got, want)
			}
		})
	}
}

type fakeProvider struct {
	message string
	err     error
	prompt  string
}

func (p *fakeProvider) Generate(_ context.Context, prompt string) (string, error) {
	p.prompt = prompt
	return p.message, p.err
}
