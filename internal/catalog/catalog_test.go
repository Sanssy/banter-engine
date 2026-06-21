package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DSanoussy/banter-engine/internal/opportunities"
)

func TestLoadCatalogLoadsDefinitiveCatalog(t *testing.T) {
	opportunityCatalog, err := LoadCatalog(filepath.Join("..", "..", "resources", "opportunities.json"))
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}
	if opportunityCatalog.Len() != 113 {
		t.Fatalf("catalog contains %d definitions, want 113", opportunityCatalog.Len())
	}
	expectedByCategory := map[string]int{
		"Ranking":     21,
		"Predictions": 26,
		"Crowd":       20,
		"MatchEvents": 26,
		"Narratives":  20,
	}
	for category, expected := range expectedByCategory {
		if definitions := opportunityCatalog.FindByCategory(category); len(definitions) != expected {
			t.Fatalf("FindByCategory(%q) returned %d definitions, want %d", category, len(definitions), expected)
		}
	}
}

func TestCatalogContainsEveryRegisteredOpportunity(t *testing.T) {
	opportunityCatalog, err := LoadCatalog(filepath.Join("..", "..", "resources", "opportunities.json"))
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}
	for _, opportunityType := range opportunities.RegisteredTypes() {
		if _, found := opportunityCatalog.FindByID(opportunityType); !found {
			t.Errorf("registered opportunity %q is missing from catalog", opportunityType)
		}
	}
}

func TestCatalogQueries(t *testing.T) {
	opportunityCatalog := loadFixture(t, validCatalog)

	definition, found := opportunityCatalog.FindByID("RankingOvertake")
	if !found || definition.Name != "Ranking Overtake" {
		t.Fatalf("FindByID() = %+v, %v", definition, found)
	}
	if definitions := opportunityCatalog.FindByCategory("Ranking"); len(definitions) != 2 {
		t.Fatalf("FindByCategory() returned %d definitions", len(definitions))
	}
	related := opportunityCatalog.FindRelated("RankingOvertake")
	if len(related) != 1 || related[0].ID != "DoubleOvertake" {
		t.Fatalf("FindRelated() = %+v", related)
	}
}

func TestLoadCatalogRejectsInvalidDefinition(t *testing.T) {
	invalid := strings.Replace(validCatalog, `"severity": 2`, `"severity": 6`, 1)
	path := writeCatalog(t, invalid)
	if _, err := LoadCatalog(path); err == nil || !strings.Contains(err.Error(), "severity") {
		t.Fatalf("LoadCatalog() error = %v, want severity error", err)
	}
}

func TestLoadCatalogRejectsDuplicateIDs(t *testing.T) {
	duplicate := strings.Replace(validCatalog, `"id": "DoubleOvertake"`, `"id": "RankingOvertake"`, 1)
	path := writeCatalog(t, duplicate)
	if _, err := LoadCatalog(path); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("LoadCatalog() error = %v, want duplicate error", err)
	}
}

func TestLoadCatalogRejectsMissingFields(t *testing.T) {
	missing := strings.Replace(validCatalog, `"tags": ["ranking"]`, `"tags": []`, 1)
	path := writeCatalog(t, missing)
	if _, err := LoadCatalog(path); err == nil || !strings.Contains(err.Error(), "tags") {
		t.Fatalf("LoadCatalog() error = %v, want missing tags error", err)
	}
}

func TestLoadCatalogRejectsUnknownRelatedOpportunity(t *testing.T) {
	unknown := strings.Replace(validCatalog, `"DoubleOvertake"]`, `"Unknown"]`, 1)
	path := writeCatalog(t, unknown)
	if _, err := LoadCatalog(path); err == nil || !strings.Contains(err.Error(), "references unknown") {
		t.Fatalf("LoadCatalog() error = %v, want unknown related opportunity error", err)
	}
}

func loadFixture(t *testing.T, data string) *Catalog {
	t.Helper()
	opportunityCatalog, err := LoadCatalog(writeCatalog(t, data))
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}
	return opportunityCatalog
}

func writeCatalog(t *testing.T, data string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "opportunities.json")
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write catalog: %v", err)
	}
	return path
}

const validCatalog = `[
  {
    "id": "RankingOvertake",
    "name": "Ranking Overtake",
    "category": "Ranking",
    "severity": 2,
    "description": "A player overtakes another player.",
    "requiredData": ["standings"],
    "trigger": {"rankImprovement": 1},
    "banterAngles": ["superiority"],
    "relatedOpportunities": ["DoubleOvertake"],
    "tags": ["ranking"]
  },
  {
    "id": "DoubleOvertake",
    "name": "Double Overtake",
    "category": "Ranking",
    "severity": 3,
    "description": "A player overtakes two players.",
    "requiredData": ["standings"],
    "trigger": {"rankImprovement": 2},
    "banterAngles": ["momentum"],
    "relatedOpportunities": [],
    "tags": ["ranking"]
  }
]`
