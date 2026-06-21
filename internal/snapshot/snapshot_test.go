package snapshot

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Sanssy/banter-engine/internal/forecasts"
	"github.com/Sanssy/banter-engine/internal/matches"
	"github.com/Sanssy/banter-engine/internal/model"
	"github.com/Sanssy/banter-engine/internal/rivalries"
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

func TestSaveAndLoadMatches(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "matches.json")
	want := []matches.Match{
		{
			MatchID: "match-1",
			Status:  "firstHalf",
			Score:   matches.Score{Home: 1},
			Events:  []matches.Event{{ID: "event-1", Type: "goal"}},
		},
	}

	if err := SaveMatches(path, want); err != nil {
		t.Fatalf("SaveMatches() error = %v", err)
	}
	got, err := LoadMatches(path)
	if err != nil {
		t.Fatalf("LoadMatches() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadMatches() = %#v, want %#v", got, want)
	}
}

func TestSaveAndLoadRivalries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "rivalries.json")
	want := []rivalries.Rivalry{{PlayerAID: "a", PlayerBID: "b", WinsA: 3, WinsB: 1}}

	if err := SaveRivalries(path, want); err != nil {
		t.Fatalf("SaveRivalries() error = %v", err)
	}
	got, err := LoadRivalries(path)
	if err != nil {
		t.Fatalf("LoadRivalries() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadRivalries() = %#v, want %#v", got, want)
	}
}

func TestSaveAndLoadForecasts(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "forecasts.json")
	want := []forecasts.Forecast{{UserID: "user-1", MatchID: "match-1", Points: 5}}

	if err := SaveForecasts(path, want); err != nil {
		t.Fatalf("SaveForecasts() error = %v", err)
	}
	got, err := LoadForecasts(path)
	if err != nil {
		t.Fatalf("LoadForecasts() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadForecasts() = %#v, want %#v", got, want)
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
