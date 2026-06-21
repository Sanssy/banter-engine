package banter

import (
	"testing"

	"github.com/DSanoussy/banter-engine/internal/catalog"
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

func TestGenerateWithDefinitionUsesCatalogDescriptionForFutureType(t *testing.T) {
	opportunity := opportunities.Opportunity{Type: "FutureOpportunity", Actor: "Sanssy", Target: "William"}
	definition := catalog.OpportunityDefinition{
		ID:          opportunity.Type,
		Description: "A future catalog-driven opportunity.",
	}
	want := "A future catalog-driven opportunity: Sanssy -> William"

	if got := GenerateWithDefinition(opportunity, definition); got != want {
		t.Fatalf("GenerateWithDefinition() = %q, want %q", got, want)
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
			name: "against the crowd",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.AgainstTheCrowd,
				Actor:  "Julien",
				Target: "France - Espagne",
			},
			want: "🎯 Julien a défié la foule et avait raison sur France - Espagne.",
		},
		{
			name: "crowd favorite",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.CrowdFavorite,
				Actor:  "France",
				Target: "70%",
			},
			want: "👥 France concentre 70% des pronostics.",
		},
		{
			name: "crowd trap",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.CrowdTrap,
				Actor:  "Espagne",
				Target: "France",
			},
			want: "🪤 La foule choisit Espagne plutôt que le favori France.",
		},
		{
			name: "popular mistake",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.PopularMistake,
				Actor:  "France",
				Target: "France - Espagne",
			},
			want: "😬 Le choix populaire France s'effondre sur France - Espagne.",
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
			name: "added time disaster",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.AddedTimeDisaster,
				Actor:  "Julien",
				Target: "France - Espagne",
			},
			want: "⏱️ Le temps additionnel coûte cher à Julien sur France - Espagne.",
		},
		{
			name: "last minute hero",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.LastMinuteHero,
				Actor:  "Sanssy",
				Target: "France - Espagne",
			},
			want: "🦸 Sanssy devient le héros de la dernière minute sur France - Espagne.",
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
		{
			name: "nemesis",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.Nemesis,
				Actor:  "Sanssy",
				Target: "William",
			},
			want: "👹 Sanssy devient la Némésis officielle de William.",
		},
		{
			name: "revenge",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.Revenge,
				Actor:  "Sanssy",
				Target: "William",
			},
			want: "🗡️ Sanssy prend sa revanche sur William.",
		},
		{
			name: "dominance",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.Dominance,
				Actor:  "Sanssy",
				Target: "William",
			},
			want: "👑 Sanssy domine désormais sa rivalité avec William.",
		},
		{
			name: "comeback season",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.ComebackSeason,
				Actor:  "Sanssy",
				Target: "3",
			},
			want: "🚀 Sanssy remonte de 3 places au classement.",
		},
		{
			name: "free fall",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.FreeFall,
				Actor:  "William",
				Target: "4",
			},
			want: "🪂 William dégringole de 4 places au classement.",
		},
		{
			name: "runaway leader",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.RunawayLeader,
				Actor:  "Amine",
				Target: "Sanssy",
			},
			want: "🏃 Amine creuse l'écart devant Sanssy.",
		},
		{
			name: "podium fight",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.PodiumFight,
				Actor:  "Sanssy",
				Target: "William",
			},
			want: "🥉 Sanssy met la pression sur William pour le podium.",
		},
		{
			name: "goal swing",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.GoalSwing,
				Actor:  "France",
				Target: "France - Espagne (1-0)",
			},
			want: "⚽ France fait basculer le score sur France - Espagne (1-0).",
		},
		{
			name: "match turnaround",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.MatchTurnaround,
				Actor:  "France",
				Target: "Espagne",
			},
			want: "🔄 France renverse complètement Espagne.",
		},
		{
			name: "equalizer chaos",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.EqualizerChaos,
				Actor:  "Espagne",
				Target: "France - Espagne",
			},
			want: "🌪️ Espagne égalise sur France - Espagne.",
		},
		{
			name: "biggest winner",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.BiggestWinner,
				Actor:  "Julien",
				Target: "+7",
			},
			want: "💰 Julien signe le plus gros gain avec +7 points.",
		},
		{
			name: "biggest loser",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.BiggestLoser,
				Actor:  "William",
				Target: "-5",
			},
			want: "📉 William subit la plus grosse perte avec -5 points.",
		},
		{
			name: "point explosion",
			opportunity: opportunities.Opportunity{
				Type:   opportunities.PointExplosion,
				Actor:  "Julien",
				Target: "+7",
			},
			want: "💥 Explosion de points pour Julien : +7.",
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
