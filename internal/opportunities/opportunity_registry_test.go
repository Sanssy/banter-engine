package opportunities

import "testing"

func TestRegisteredTypesAreUnique(t *testing.T) {
	seen := make(map[string]struct{}, len(registeredTypes))
	for _, opportunityType := range RegisteredTypes() {
		if opportunityType == "" {
			t.Fatal("RegisteredTypes() contains an empty type")
		}
		if _, duplicate := seen[opportunityType]; duplicate {
			t.Errorf("RegisteredTypes() contains duplicate %q", opportunityType)
		}
		seen[opportunityType] = struct{}{}
	}
}
