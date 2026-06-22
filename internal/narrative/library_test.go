package narrative

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLibraryAndSelectCompatibleExamples(t *testing.T) {
	path := filepath.Join(t.TempDir(), "examples.json")
	data := `{
		"version": 2,
		"examples": [
			{"category":"Ranking","angle":"Rise","facts":"A gagne trois places.","message":"A prend l'ascenseur."},
			{"category":"Crowd","angle":"Rise","facts":"B progresse.","message":"B monte."},
			{"category":"Ranking","angle":"Rise","facts":"C entre sur le podium.","message":"C voit le podium."},
			{"category":"Crowd","angle":"CrowdWrong","facts":"La foule se trompe.","message":"Le consensus tombe."}
		]
	}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}

	library, err := LoadLibrary(path)
	if err != nil {
		t.Fatalf("LoadLibrary() error = %v", err)
	}
	if library.Version() != 2 || library.Len() != 4 {
		t.Fatalf("loaded library version=%d len=%d", library.Version(), library.Len())
	}

	selected := library.Select("Ranking", Rise, 3)
	if len(selected) != 3 {
		t.Fatalf("Select() returned %d examples, want 3", len(selected))
	}
	if selected[0].Category != "Ranking" || selected[1].Category != "Ranking" {
		t.Fatalf("Select() did not prioritize category: %#v", selected)
	}
	if selected[2].Angle != Rise {
		t.Fatalf("Select() returned incompatible angle: %#v", selected[2])
	}
}

func TestSelectNeverReturnsMoreThanFiveExamples(t *testing.T) {
	library := &Library{examples: []Example{
		{Category: "Crowd", Angle: CrowdWrong, Facts: "1", Message: "1"},
		{Category: "Crowd", Angle: CrowdWrong, Facts: "2", Message: "2"},
		{Category: "Crowd", Angle: CrowdWrong, Facts: "3", Message: "3"},
		{Category: "Crowd", Angle: CrowdWrong, Facts: "4", Message: "4"},
		{Category: "Crowd", Angle: CrowdWrong, Facts: "5", Message: "5"},
		{Category: "Crowd", Angle: CrowdWrong, Facts: "6", Message: "6"},
	}}

	if got := len(library.Select("Crowd", CrowdWrong, 100)); got != MaxExamples {
		t.Fatalf("Select() returned %d examples, want maximum %d", got, MaxExamples)
	}
}

func TestRepositoryNarrativeLibraryIsValid(t *testing.T) {
	library, err := LoadLibrary(filepath.Join("..", "..", "resources", "narratives", "examples.json"))
	if err != nil {
		t.Fatalf("repository narrative library is invalid: %v", err)
	}
	if library.Len() < 18 {
		t.Fatalf("repository narrative library has %d examples, want at least 18", library.Len())
	}
	for _, angle := range []Angle{
		CrowdWrong,
		MinorityVictory,
		FallFromGrace,
		Rise,
		Curse,
		Dominance,
	} {
		if got := len(library.Select("", angle, MaxExamples)); got != MaxExamples {
			t.Errorf("repository library provides %d examples for %s, want %d", got, angle, MaxExamples)
		}
	}
}
