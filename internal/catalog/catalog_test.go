package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DSanoussy/banter-engine/internal/opportunities"
)

func TestLoadOpportunityCatalogRegistersEverySupportedOpportunity(t *testing.T) {
	definitions, err := LoadOpportunityCatalog(filepath.Join("..", "..", "resources", "opportunities.json"))
	if err != nil {
		t.Fatalf("LoadOpportunityCatalog() error = %v", err)
	}
	expected := []string{
		opportunities.RankingOvertake,
		opportunities.EnteredTop3,
		opportunities.ExitedTop3,
		opportunities.LeaderUnderPressure,
		opportunities.LastPlaceLocked,
		opportunities.ComebackSeason,
		opportunities.FreeFall,
		opportunities.RunawayLeader,
		opportunities.PodiumFight,
		opportunities.HugeUpset,
		opportunities.EveryoneWasWrong,
		opportunities.TheChosenOne,
		opportunities.AgainstTheCrowd,
		opportunities.CrowdFavorite,
		opportunities.CrowdTrap,
		opportunities.PopularMistake,
		opportunities.PredictionMassacre,
		opportunities.HotStreak,
		opportunities.ColdStreak,
		opportunities.MatchStarted,
		opportunities.MatchEnded,
		opportunities.ScoreChanged,
		opportunities.ImportantMatchEvent,
		opportunities.GoalSwing,
		opportunities.MatchTurnaround,
		opportunities.EqualizerChaos,
		opportunities.BiggestWinner,
		opportunities.BiggestLoser,
		opportunities.PointExplosion,
		opportunities.NinetiethMinuteHeartbreak,
		opportunities.AddedTimeDisaster,
		opportunities.LastMinuteHero,
		opportunities.VARVictim,
		opportunities.RedCardDisaster,
		opportunities.Nemesis,
		opportunities.Revenge,
		opportunities.Dominance,
	}
	if len(definitions) != len(expected) {
		t.Fatalf("catalog contains %d definitions, want %d", len(definitions), len(expected))
	}
	for _, id := range expected {
		definition, found := FindOpportunity(definitions, id)
		if !found || definition.ID != id {
			t.Fatalf("opportunity %q is missing from catalog", id)
		}
	}
}

func TestLoadOpportunityCatalogRejectsDuplicates(t *testing.T) {
	path := filepath.Join(t.TempDir(), "opportunities.json")
	data := `[
		{"id":"Duplicate","category":"Test","severity":1,"description":"One","tags":[]},
		{"id":"Duplicate","category":"Test","severity":1,"description":"Two","tags":[]}
	]`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write catalog: %v", err)
	}

	if _, err := LoadOpportunityCatalog(path); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("LoadOpportunityCatalog() error = %v, want duplicate error", err)
	}
}

func TestValidateOpportunityRejectsUnknownType(t *testing.T) {
	definitions := []OpportunityDefinition{{ID: opportunities.RankingOvertake}}
	_, err := ValidateOpportunity(definitions, opportunities.Opportunity{Type: "Unknown"})
	if err == nil || !strings.Contains(err.Error(), "unknown opportunity") {
		t.Fatalf("ValidateOpportunity() error = %v", err)
	}
}
