package password

import (
	"database/sql"
	"path/filepath"
	"time"

	"pocketunzip/internal/security"
)

type PasswordRecord struct {
	ID                int64
	ArchivePath       string
	ArchiveName       string
	ArchiveHash       string
	EncryptedPassword []byte
	SuccessCount      int
	LastUsedAt        time.Time
	CreatedAt         time.Time
}

func Save(db *sql.DB, archivePath, password string) error {
	encrypted, err := security.Encrypt([]byte(password))
	if err != nil {
		return err
	}

	name := filepath.Base(archivePath)

	_, err = db.Exec(
		`INSERT INTO password_records (archive_path, archive_name, encrypted_password)
		 VALUES (?, ?, ?)`,
		archivePath,
		name,
		encrypted,
	)
	return err
}

func Match(db *sql.DB, archivePath string) ([]string, error) {
	name := filepath.Base(archivePath)

	rows, err := db.Query(
		`SELECT encrypted_password FROM password_records
		 WHERE archive_path = ? OR archive_name = ?
		 ORDER BY success_count DESC, last_used_at DESC
		 LIMIT 10`,
		archivePath,
		name,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passwords []string
	for rows.Next() {
		var encrypted []byte
		if err := rows.Scan(&encrypted); err != nil {
			continue
		}
		decrypted, err := security.Decrypt(encrypted)
		if err != nil {
			continue
		}
		passwords = append(passwords, string(decrypted))
	}
	return passwords, rows.Err()
}

func UpdateSuccess(db *sql.DB, archivePath, password string) error {
	name := filepath.Base(archivePath)
	encrypted, err := security.Encrypt([]byte(password))
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`UPDATE password_records
		 SET success_count = success_count + 1, last_used_at = ?
		 WHERE (archive_path = ? OR archive_name = ?) AND encrypted_password = ?`,
		time.Now(),
		archivePath,
		name,
		encrypted,
	)
	return err
}
