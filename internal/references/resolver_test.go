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
	for _, expected := range []string{
		"club_lookup id=mpp_championship_club_614 found=true name=Argentine",
		"user_lookup id=user_11291094 found=true username=LeDaveCoinCoin",
	} {
		if !strings.Contains(output.String(), expected) {
			t.Errorf("resolver log does not contain %q: %s", expected, output.String())
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
	if expected := "club_lookup id=mpp_championship_club_999 found=false name=Equipe inconnue available_keys=[mpp_championship_club_367 mpp_championship_club_614]"; !strings.Contains(output.String(), expected) {
		t.Errorf("resolver log does not contain %q: %s", expected, output.String())
	}
}
