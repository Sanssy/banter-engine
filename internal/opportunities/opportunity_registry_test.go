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

func TestEnsureIdentityBuildsStableBusinessKey(t *testing.T) {
	op := Opportunity{
		Type:    ImportantMatchEvent,
		Actor:   "France - Espagne",
		Target:  "goal 90'",
		MatchID: "match-1",
		EventID: "event-1",
	}

	first := EnsureIdentity(op)
	second := EnsureIdentity(op)
	if first.Key == "" {
		t.Fatal("EnsureIdentity() returned an empty key")
	}
	if first.Key != second.Key {
		t.Fatalf("EnsureIdentity() keys differ: %q != %q", first.Key, second.Key)
	}

	otherMatch := op
	otherMatch.MatchID = "match-2"
	if EnsureIdentity(otherMatch).Key == first.Key {
		t.Fatal("different matches produced the same opportunity key")
	}
	otherEvent := op
	otherEvent.EventID = "event-2"
	if EnsureIdentity(otherEvent).Key == first.Key {
		t.Fatal("different events produced the same opportunity key")
	}
}
