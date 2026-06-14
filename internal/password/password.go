package password

import (
	"database/sql"
	"path/filepath"
	"time"

	"pocketunzip/internal/security"
)

type PasswordRecord struct {
	ID                int64     `json:"id"`
	ArchivePath       string    `json:"archivePath"`
	ArchiveName       string    `json:"archiveName"`
	ArchiveHash       string    `json:"archiveHash"`
	EncryptedPassword []byte    `json:"-"`
	SuccessCount      int       `json:"successCount"`
	LastUsedAt        time.Time `json:"lastUsedAt"`
	CreatedAt         time.Time `json:"createdAt"`
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

	// 获取所有匹配的记录，解密后比较密码
	rows, err := db.Query(
		`SELECT id, encrypted_password FROM password_records
		 WHERE archive_path = ? OR archive_name = ?`,
		archivePath,
		name,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	var recordID int64
	for rows.Next() {
		var id int64
		var encrypted []byte
		if err := rows.Scan(&id, &encrypted); err != nil {
			continue
		}
		decrypted, err := security.Decrypt(encrypted)
		if err != nil {
			continue
		}
		if string(decrypted) == password {
			recordID = id
			break
		}
	}

	if recordID == 0 {
		return nil // 没有找到匹配的记录
	}

	_, err = db.Exec(
		`UPDATE password_records
		 SET success_count = success_count + 1, last_used_at = ?
		 WHERE id = ?`,
		time.Now(),
		recordID,
	)
	return err
}

// ListAll 获取所有密码记录（密码库管理）
func ListAll(db *sql.DB) ([]PasswordRecord, error) {
	rows, err := db.Query(
		`SELECT id, archive_path, archive_name, archive_hash, encrypted_password, success_count, last_used_at, created_at
		 FROM password_records
		 ORDER BY success_count DESC, last_used_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []PasswordRecord
	for rows.Next() {
		var r PasswordRecord
		var lastUsedAt, createdAt sql.NullTime
		err := rows.Scan(&r.ID, &r.ArchivePath, &r.ArchiveName, &r.ArchiveHash, &r.EncryptedPassword, &r.SuccessCount, &lastUsedAt, &createdAt)
		if err != nil {
			continue
		}
		if lastUsedAt.Valid {
			r.LastUsedAt = lastUsedAt.Time
		}
		if createdAt.Valid {
			r.CreatedAt = createdAt.Time
		}
		records = append(records, r)
	}
	return records, rows.Err()
}

// DeleteByID 删除密码记录
func DeleteByID(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM password_records WHERE id = ?", id)
	return err
}

// UpdatePath 更新密码记录的关联路径
func UpdatePath(db *sql.DB, id int64, newArchivePath string) error {
	name := filepath.Base(newArchivePath)
	_, err := db.Exec(
		"UPDATE password_records SET archive_path = ?, archive_name = ? WHERE id = ?",
		newArchivePath, name, id,
	)
	return err
}

// PasswordStats 密码统计信息
type PasswordStats struct {
	TotalRecords int `json:"totalRecords"`
	TotalUsed    int `json:"totalUsed"`
}

// GetStats 获取密码使用统计
func GetStats(db *sql.DB) (PasswordStats, error) {
	var stats PasswordStats
	err := db.QueryRow("SELECT COUNT(*), COALESCE(SUM(success_count), 0) FROM password_records").Scan(&stats.TotalRecords, &stats.TotalUsed)
	return stats, err
}
