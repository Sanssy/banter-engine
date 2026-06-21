package catalog

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type OpportunityDefinition struct {
	ID                   string         `json:"id"`
	Name                 string         `json:"name"`
	Category             string         `json:"category"`
	Severity             int            `json:"severity"`
	Description          string         `json:"description"`
	RequiredData         []string       `json:"requiredData"`
	Trigger              map[string]any `json:"trigger"`
	BanterAngles         []string       `json:"banterAngles"`
	RelatedOpportunities []string       `json:"relatedOpportunities"`
	Tags                 []string       `json:"tags"`
}

type Catalog struct {
	definitions []OpportunityDefinition
	byID        map[string]OpportunityDefinition
}

func LoadCatalog(path string) (*Catalog, error) {
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

	byID := make(map[string]OpportunityDefinition, len(definitions))
	for _, definition := range definitions {
		byID[definition.ID] = definition
	}
	return &Catalog{definitions: definitions, byID: byID}, nil
}

func (c *Catalog) Len() int {
	return len(c.definitions)
}

func (c *Catalog) FindByID(id string) (OpportunityDefinition, bool) {
	definition, found := c.byID[id]
	return definition, found
}

func (c *Catalog) FindByCategory(category string) []OpportunityDefinition {
	var matches []OpportunityDefinition
	for _, definition := range c.definitions {
		if definition.Category == category {
			matches = append(matches, definition)
		}
	}
	return matches
}

func (c *Catalog) FindRelated(id string) []OpportunityDefinition {
	definition, found := c.FindByID(id)
	if !found {
		return nil
	}

	result := make([]OpportunityDefinition, 0, len(definition.RelatedOpportunities))
	for _, relatedID := range definition.RelatedOpportunities {
		if related, found := c.FindByID(relatedID); found {
			result = append(result, related)
		}
	}
	return result
}

func validateDefinitions(definitions []OpportunityDefinition) error {
	if len(definitions) == 0 {
		return fmt.Errorf("opportunity catalog is empty")
	}
	seen := make(map[string]struct{}, len(definitions))
	for i, definition := range definitions {
		if strings.TrimSpace(definition.ID) == "" || strings.TrimSpace(definition.Name) == "" {
			return fmt.Errorf("opportunity definition at index %d is missing id or name", i)
		}
		if !validCategory(definition.Category) {
			return fmt.Errorf("opportunity %q has invalid category %q", definition.ID, definition.Category)
		}
		if definition.Severity < 1 || definition.Severity > 5 {
			return fmt.Errorf("opportunity %q has invalid severity %d", definition.ID, definition.Severity)
		}
		if strings.TrimSpace(definition.Description) == "" {
			return fmt.Errorf("opportunity %q has no description", definition.ID)
		}
		if len(definition.RequiredData) == 0 || len(definition.Trigger) == 0 || len(definition.BanterAngles) == 0 {
			return fmt.Errorf("opportunity %q has incomplete detection metadata", definition.ID)
		}
		if definition.RelatedOpportunities == nil {
			return fmt.Errorf("opportunity %q is missing related opportunities", definition.ID)
		}
		if len(definition.Tags) == 0 {
			return fmt.Errorf("opportunity %q has no tags", definition.ID)
		}
		if _, duplicate := seen[definition.ID]; duplicate {
			return fmt.Errorf("duplicate opportunity %q", definition.ID)
		}
		seen[definition.ID] = struct{}{}
	}
	return nil
}

func validCategory(category string) bool {
	switch category {
	case "Ranking", "Predictions", "Crowd", "MatchEvents", "Narratives":
		return true
	default:
		return false
	}
}
