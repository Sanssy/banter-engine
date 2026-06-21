package context

import (
	"encoding/json"

	"github.com/DSanoussy/banter-engine/internal/catalog"
	"github.com/DSanoussy/banter-engine/internal/opportunities"
)

func BuildLLMContext(
	opportunity opportunities.Opportunity,
	definition catalog.OpportunityDefinition,
) string {
	data, err := json.Marshal(struct {
		ID          string   `json:"id"`
		Category    string   `json:"category"`
		Severity    int      `json:"severity"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
		Actor       string   `json:"actor"`
		Target      string   `json:"target"`
	}{
		ID:          definition.ID,
		Category:    definition.Category,
		Severity:    definition.Severity,
		Description: definition.Description,
		Tags:        definition.Tags,
		Actor:       opportunity.Actor,
		Target:      opportunity.Target,
	})
	if err != nil {
		return "{}"
	}
	return string(data)
}
