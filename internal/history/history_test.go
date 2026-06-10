package history

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE extract_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		archive_path TEXT NOT NULL,
		output_dir TEXT NOT NULL,
		success INTEGER NOT NULL,
		used_password INTEGER DEFAULT 0,
		error_message TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestRecord(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	h := ExtractHistory{
		ArchivePath:  "test.zip",
		OutputDir:    "output",
		Success:      true,
		UsedPassword: false,
	}

	err := Record(db, h)
	if err != nil {
		t.Fatalf("Record() error = %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM extract_history").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 record, got %d", count)
	}
}

func TestList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// 插入测试数据
	for i := 0; i < 5; i++ {
		Record(db, ExtractHistory{
			ArchivePath: "test" + string(rune('0'+i)) + ".zip",
			OutputDir:   "output",
			Success:     true,
		})
	}

	histories, err := List(db, 3)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(histories) != 3 {
		t.Errorf("expected 3 records, got %d", len(histories))
	}
}
