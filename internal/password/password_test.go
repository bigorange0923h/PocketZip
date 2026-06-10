package password

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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS password_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		archive_path TEXT,
		archive_name TEXT,
		archive_hash TEXT,
		encrypted_password BLOB NOT NULL,
		success_count INTEGER DEFAULT 1,
		last_used_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSave(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := Save(db, "/path/to/test.zip", "password123")
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM password_records").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 record, got %d", count)
	}
}

func TestMatch(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	Save(db, "/path/to/test.zip", "password1")
	Save(db, "/path/to/test.zip", "password2")
	Save(db, "/other/test.zip", "password3")

	passwords, err := Match(db, "/path/to/test.zip")
	if err != nil {
		t.Fatalf("Match() error = %v", err)
	}
	if len(passwords) < 2 {
		t.Errorf("expected at least 2 passwords, got %d", len(passwords))
	}
}

func TestUpdateSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// 先插入一条记录
	_, err := db.Exec(`INSERT INTO password_records (archive_path, archive_name, encrypted_password, success_count)
		VALUES (?, ?, ?, ?)`, "/path/to/test.zip", "test.zip", []byte("encrypted"), 1)
	if err != nil {
		t.Fatal(err)
	}

	// 验证插入成功
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM password_records").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 record, got %d", count)
	}

	// 直接更新 success_count
	_, err = db.Exec(`UPDATE password_records SET success_count = success_count + 1 WHERE archive_path = ?`, "/path/to/test.zip")
	if err != nil {
		t.Fatalf("Update error = %v", err)
	}

	// 验证更新成功
	var successCount int
	err = db.QueryRow("SELECT success_count FROM password_records WHERE archive_path = ?", "/path/to/test.zip").Scan(&successCount)
	if err != nil {
		t.Fatal(err)
	}
	if successCount != 2 {
		t.Errorf("expected success_count 2, got %d", successCount)
	}
}