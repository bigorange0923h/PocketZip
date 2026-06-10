package history

import (
	"database/sql"
	"time"
)

type ExtractHistory struct {
	ID           int64
	ArchivePath  string
	OutputDir    string
	Success      bool
	UsedPassword bool
	ErrorMessage string
	CreatedAt    time.Time
}

func Record(db *sql.DB, h ExtractHistory) error {
	_, err := db.Exec(
		`INSERT INTO extract_history (archive_path, output_dir, success, used_password, error_message)
		 VALUES (?, ?, ?, ?, ?)`,
		h.ArchivePath,
		h.OutputDir,
		boolToInt(h.Success),
		boolToInt(h.UsedPassword),
		h.ErrorMessage,
	)
	return err
}

func List(db *sql.DB, limit int) ([]ExtractHistory, error) {
	rows, err := db.Query(
		`SELECT id, archive_path, output_dir, success, used_password, error_message, created_at
		 FROM extract_history
		 ORDER BY created_at DESC
		 LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []ExtractHistory
	for rows.Next() {
		var h ExtractHistory
		var success, usedPassword int
		err := rows.Scan(&h.ID, &h.ArchivePath, &h.OutputDir, &success, &usedPassword, &h.ErrorMessage, &h.CreatedAt)
		if err != nil {
			return nil, err
		}
		h.Success = success != 0
		h.UsedPassword = usedPassword != 0
		histories = append(histories, h)
	}
	return histories, rows.Err()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
