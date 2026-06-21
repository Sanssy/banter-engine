package banter

import (
	"testing"

	"github.com/DSanoussy/banter-engine/internal/opportunities"
)

func TestGenerateRankingOvertake(t *testing.T) {
	opportunity := opportunities.Opportunity{
		Type:   opportunities.RankingOvertake,
		Actor:  "Sanssy",
		Target: "William",
	}
	want := "📈 Sanssy dépasse William.\n\nWilliam aperçoit désormais son dos."

	if got := Generate(opportunity); got != want {
		t.Fatalf("Generate() = %q, want %q", got, want)
	}
}

func TestGenerateUnknownOpportunity(t *testing.T) {
	opportunity := opportunities.Opportunity{
		Type:   "Unknown",
		Actor:  "Sanssy",
		Target: "William",
	}
	want := "Unknown: Sanssy -> William"

	if got := Generate(opportunity); got != want {
		t.Fatalf("Generate() = %q, want %q", got, want)
	}
}
