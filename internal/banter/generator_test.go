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

func TestGenerateSupportedOpportunities(t *testing.T) {
	tests := []struct {
		name        string
		opportunity opportunities.Opportunity
		want        string
	}{
		{
			name:        "entered top 3",
			opportunity: opportunities.Opportunity{Type: opportunities.EnteredTop3, Actor: "Sanssy"},
			want:        "🏆 Sanssy entre dans le top 3.",
		},
		{
			name:        "exited top 3",
			opportunity: opportunities.Opportunity{Type: opportunities.ExitedTop3, Actor: "William"},
			want:        "📉 William quitte le top 3.",
		},
		{
			name: "leader under pressure",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.LeaderUnderPressure,
				Actor:  "Amine",
				Target: "Sanssy",
			},
			want: "🔥 Amine voit Sanssy revenir dans son rétroviseur.",
		},
		{
			name:        "last place locked",
			opportunity: opportunities.Opportunity{Type: opportunities.LastPlaceLocked, Actor: "Julien"},
			want:        "🔒 Julien conserve solidement la dernière place.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Generate(tt.opportunity); got != tt.want {
				t.Fatalf("Generate() = %q, want %q", got, tt.want)
			}
		})
	}
}
