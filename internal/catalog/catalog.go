package catalog

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/DSanoussy/banter-engine/internal/opportunities"
)

type OpportunityDefinition struct {
	ID          string   `json:"id"`
	Category    string   `json:"category"`
	Severity    int      `json:"severity"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func LoadOpportunityCatalog(path string) ([]OpportunityDefinition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read opportunity catalog: %w", err)
	}

	var definitions []OpportunityDefinition
	if err := json.Unmarshal(data, &definitions); err != nil {
		return nil, fmt.Errorf("decode opportunity catalog: %w", err)
	}
	if err := validateDefinitions(definitions); err != nil {
		return nil, err
	}
	return definitions, nil
}

func FindOpportunity(
	definitions []OpportunityDefinition,
	id string,
) (OpportunityDefinition, bool) {
	for _, definition := range definitions {
		if definition.ID == id {
			return definition, true
		}
	}
	return OpportunityDefinition{}, false
}

func ValidateOpportunity(
	definitions []OpportunityDefinition,
	opportunity opportunities.Opportunity,
) (OpportunityDefinition, error) {
	definition, found := FindOpportunity(definitions, opportunity.Type)
	if !found {
		return OpportunityDefinition{}, fmt.Errorf("unknown opportunity %q", opportunity.Type)
	}
	return definition, nil
}

func validateDefinitions(definitions []OpportunityDefinition) error {
	if len(definitions) == 0 {
		return fmt.Errorf("opportunity catalog is empty")
	}
	seen := make(map[string]struct{}, len(definitions))
	for i, definition := range definitions {
		if definition.ID == "" || definition.Category == "" || definition.Description == "" {
			return fmt.Errorf("opportunity definition at index %d is incomplete", i)
		}
		if definition.Severity < 1 || definition.Severity > 5 {
			return fmt.Errorf("opportunity %q has invalid severity %d", definition.ID, definition.Severity)
		}
		if _, duplicate := seen[definition.ID]; duplicate {
			return fmt.Errorf("duplicate opportunity %q", definition.ID)
		}
		seen[definition.ID] = struct{}{}
	}
	return nil
}
