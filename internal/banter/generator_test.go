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
		{
			name: "huge upset",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.HugeUpset,
				Actor:  "Espagne",
				Target: "France",
			},
			want: "⚠️ Espagne fait tomber le favori France.",
		},
		{
			name:        "everyone was wrong",
			opportunity: opportunities.Opportunity{Type: opportunities.EveryoneWasWrong, Actor: "France - Espagne"},
			want:        "📉 France - Espagne : l'intelligence collective a pris un coup.",
		},
		{
			name: "the chosen one",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.TheChosenOne,
				Actor:  "Julien",
				Target: "France - Espagne",
			},
			want: "🔮 Julien était le seul à y croire sur France - Espagne.",
		},
		{
			name:        "prediction massacre",
			opportunity: opportunities.Opportunity{Type: opportunities.PredictionMassacre, Actor: "France - Espagne"},
			want:        "☠️ Extinction des pronostics détectée sur France - Espagne.",
		},
		{
			name: "hot streak",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.HotStreak,
				Actor:  "Julien",
				Target: "5",
			},
			want: "🔥 Julien enchaîne 5 pronostics réussis.",
		},
		{
			name: "cold streak",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.ColdStreak,
				Actor:  "William",
				Target: "6",
			},
			want: "🥶 William reste sur 6 échecs consécutifs.",
		},
		{
			name:        "match started",
			opportunity: opportunities.Opportunity{Type: opportunities.MatchStarted, Actor: "France - Espagne"},
			want:        "⚽ Coup d'envoi pour France - Espagne.",
		},
		{
			name:        "match ended",
			opportunity: opportunities.Opportunity{Type: opportunities.MatchEnded, Actor: "France - Espagne"},
			want:        "🏁 Fin du match France - Espagne.",
		},
		{
			name: "score changed",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.ScoreChanged,
				Actor:  "France - Espagne",
				Target: "1-0",
			},
			want: "⚽ France - Espagne : le score passe à 1-0.",
		},
		{
			name: "important match event",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.ImportantMatchEvent,
				Actor:  "France - Espagne",
				Target: "goal 42'",
			},
			want: "🚨 France - Espagne : goal 42'.",
		},
		{
			name: "ninetieth minute heartbreak",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.NinetiethMinuteHeartbreak,
				Actor:  "Julien",
				Target: "France - Espagne",
			},
			want: "💔 Julien perd son prono parfait dans les derniers instants de France - Espagne.",
		},
		{
			name: "var victim",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.VARVictim,
				Actor:  "Julien",
				Target: "France - Espagne",
			},
			want: "📺 La VAR brise le pronostic de Julien sur France - Espagne.",
		},
		{
			name: "red card disaster",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.RedCardDisaster,
				Actor:  "Julien",
				Target: "France - Espagne",
			},
			want: "🟥 Le rouge met le pronostic de Julien en danger sur France - Espagne.",
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
