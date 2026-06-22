package narrative

import (
	"testing"

	"github.com/Sanssy/banter-engine/internal/opportunities"
)

func TestEveryRegisteredOpportunityHasNarrativeAngle(t *testing.T) {
	for _, opportunityType := range opportunities.RegisteredTypes() {
		angle := ForOpportunity(opportunityType)
		if angle == "" {
			t.Errorf("ForOpportunity(%q) returned no angle", opportunityType)
		}
		if angle.Guidance() == "" {
			t.Errorf("angle %q has no prompt guidance", angle)
		}
	}
}

func TestForOpportunityUsesExpectedCoreAngles(t *testing.T) {
	tests := map[string]Angle{
		opportunities.EveryoneWasWrong: CrowdWrong,
		opportunities.TheChosenOne:     MinorityVictory,
		opportunities.HugeUpset:        FallFromGrace,
		opportunities.RankingOvertake:  Rise,
		opportunities.ColdStreak:       Curse,
		opportunities.Dominance:        Dominance,
	}

	for opportunityType, want := range tests {
		if got := ForOpportunity(opportunityType); got != want {
			t.Errorf("ForOpportunity(%q) = %q, want %q", opportunityType, got, want)
		}
	}
}
