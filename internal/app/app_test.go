package app

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"pocketunzip/internal/db"
	"pocketunzip/internal/history"
	"pocketunzip/internal/password"
)

func TestMain(m *testing.M) {
	if os.Getenv("POCKETUNZIP_FAKE_7Z") == "1" {
		runFake7z()
		return
	}

	os.Exit(m.Run())
}

func runFake7z() {
	for _, arg := range os.Args[1:] {
		if arg == "-psecret" {
			fmt.Println("Everything is Ok")
			os.Exit(0)
		}
	}

	fmt.Fprintln(os.Stderr, "ERROR: Wrong password")
	os.Exit(2)
}

func TestExtractTriesSavedPasswordCandidates(t *testing.T) {
	t.Setenv("POCKETUNZIP_FAKE_7Z", "1")

	database, err := db.Init(filepath.Join(t.TempDir(), "pocketunzip-test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	archivePath := filepath.Join(t.TempDir(), "secret.zip")
	outputDir := filepath.Join(t.TempDir(), "secret")
	if err := password.Save(database, archivePath, "secret"); err != nil {
		t.Fatal(err)
	}

	app := NewApp(os.Args[0], database)
	if err := app.Extract(archivePath, outputDir); err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	records, err := history.List(database, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 history record, got %d", len(records))
	}
	if !records[0].Success {
		t.Fatal("expected successful history record")
	}
	if !records[0].UsedPassword {
		t.Fatal("expected history to record saved password usage")
	}
}
