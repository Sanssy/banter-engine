package references

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Sanssy/banter-engine/internal/model"
)

func TestResolverResolvesClubsAndUsers(t *testing.T) {
	var output bytes.Buffer
	resolver := New(&output)
	resolver.RegisterClub("mpp_championship_club_614", "Argentine")
	resolver.RegisterUsers([]model.Standing{{UserID: "user_11291094", Name: "LeDaveCoinCoin"}})

	if got := resolver.ClubName("mpp_championship_club_614"); got != "Argentine" {
		t.Fatalf("ClubName() = %q, want Argentine", got)
	}
	if got := resolver.UserName("user_11291094"); got != "LeDaveCoinCoin" {
		t.Fatalf("UserName() = %q, want LeDaveCoinCoin", got)
	}
	if expected := "user reference loaded users_count=1"; !strings.Contains(output.String(), expected) {
		t.Errorf("resolver log does not contain %q: %s", expected, output.String())
	}
	for _, removed := range []string{"club_lookup", "user_lookup"} {
		if strings.Contains(output.String(), removed) {
			t.Errorf("resolver log still contains investigation marker %q: %s", removed, output.String())
		}
	}
}

func TestResolverDoesNotExposeUnknownTechnicalIDs(t *testing.T) {
	var output bytes.Buffer
	resolver := New(&output)
	resolver.RegisterClub("mpp_championship_club_614", "Argentine")
	resolver.RegisterClub("mpp_championship_club_367", "France")
	for _, technicalID := range []string{"mpp_championship_club_999", "user_999"} {
		if got := resolver.Resolve(technicalID); strings.Contains(got, technicalID) {
			t.Errorf("Resolve(%q) exposed technical ID as %q", technicalID, got)
		}
	}
	if output.Len() != 0 {
		t.Errorf("unknown reference resolution emitted investigation logs: %s", output.String())
	}
}
