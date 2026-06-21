package notify

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Sanssy/banter-engine/internal/catalog"
	"github.com/Sanssy/banter-engine/internal/opportunities"
)

var opportunityEmoji = map[string]string{
	opportunities.BiggestWinner:             "🔥",
	opportunities.PointExplosion:            "🔥",
	opportunities.HotStreak:                 "🔥",
	opportunities.TheChosenOne:              "🔥",
	opportunities.BiggestLoser:              "📉",
	opportunities.FreeFall:                  "📉",
	opportunities.ColdStreak:                "📉",
	opportunities.AddedTimeDisaster:         "📉",
	opportunities.NinetiethMinuteHeartbreak: "💔",
	opportunities.Nemesis:                   "💔",
	opportunities.RankingOvertake:           "📈",
	opportunities.EnteredTop3:               "📈",
	opportunities.ComebackSeason:            "📈",
	opportunities.HugeUpset:                 "⚠️",
	opportunities.EveryoneWasWrong:          "⚠️",
	opportunities.CrowdTrap:                 "⚠️",
	opportunities.PredictionMassacre:        "⚠️",
	opportunities.MatchTurnaround:           "⚠️",
	opportunities.GoalSwing:                 "⚽",
	opportunities.ScoreChanged:              "⚽",
	opportunities.ImportantMatchEvent:       "⚽",
	opportunities.Revenge:                   "🎯",
	opportunities.Dominance:                 "🎯",
	opportunities.LastMinuteHero:            "🎯",
}

// NightSummary builds a morning digest message from overnight opportunities.
func NightSummary(ops []opportunities.Opportunity, cat *catalog.Catalog) string {
	if len(ops) == 0 {
		return ""
	}

	selected := SelectTop(ops, cat, 6)

	var lines []string
	lines = append(lines, "**Pendant la nuit :**")
	for _, op := range selected {
		def, found := cat.FindByID(op.Type)
		if !found {
			continue
		}
		emoji := opportunityEmoji[op.Type]
		if emoji == "" {
			emoji = "•"
		}
		line := formatLine(emoji, op, def)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func formatLine(emoji string, op opportunities.Opportunity, def catalog.OpportunityDefinition) string {
	switch {
	case op.Actor != "" && op.Target != "":
		return fmt.Sprintf("%s %s vs %s — %s", emoji, op.Actor, op.Target, def.Name)
	case op.Actor != "":
		return fmt.Sprintf("%s %s — %s", emoji, op.Actor, def.Name)
	default:
		return fmt.Sprintf("%s %s", emoji, def.Name)
	}
}

// LiveDigest builds a short summary for a mid-game cycle.
func LiveDigest(ops []opportunities.Opportunity, cat *catalog.Catalog) string {
	selected := SelectTop(ops, cat, MaxNotificationsPerRun)
	if len(selected) == 0 {
		return ""
	}

	// Sort by severity descending for consistent output
	type item struct {
		op   opportunities.Opportunity
		sev  int
		line string
	}
	items := make([]item, 0, len(selected))
	for _, op := range selected {
		def, found := cat.FindByID(op.Type)
		if !found {
			continue
		}
		emoji := opportunityEmoji[op.Type]
		if emoji == "" {
			emoji = "•"
		}
		items = append(items, item{op, def.Severity, formatLine(emoji, op, def)})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].sev > items[j].sev })

	lines := make([]string, 0, len(items))
	for _, it := range items {
		lines = append(lines, it.line)
	}
	return strings.Join(lines, "\n")
}
