package snapshot

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/DSanoussy/banter-engine/internal/model"
)

func TestSaveAndLoadStandings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "standings.json")
	want := []model.Standing{
		{UserID: "user-1", Name: "Amine", Rank: 1, Points: 1676},
		{UserID: "user-2", Name: "Killian", Rank: 2, Points: 1659},
	}

	if err := SaveStandings(path, want); err != nil {
		t.Fatalf("SaveStandings() error = %v", err)
	}

	got, err := LoadStandings(path)
	if err != nil {
		t.Fatalf("LoadStandings() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadStandings() = %#v, want %#v", got, want)
	}
}

func TestLoadStandingsReturnsEmptyWhenFileDoesNotExist(t *testing.T) {
	standings, err := LoadStandings(filepath.Join(t.TempDir(), "standings.json"))
	if err != nil {
		t.Fatalf("LoadStandings() error = %v", err)
	}
	if len(standings) != 0 {
		t.Fatalf("LoadStandings() returned %d standings, want 0", len(standings))
	}
}
