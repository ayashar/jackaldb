package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

// generatelogid membuat id unik untuk setiap log entry
// kombinasi timestamp nanosecond + angka random untuk menghindari collision
func generateLogID() string {
	return fmt.Sprintf("log_%d_%d", time.Now().UnixNano(), rand.Intn(9999))
}

// logevent menyimpan satu event ke tabel logs di sqlite
// dipanggil dari http handler saat ada event masuk dari node.js
func LogEvent(db *sql.DB, event, userID, dbID, status, detail, errMsg string) error {
	query := `
		INSERT INTO logs (log_id, timestamp, event, user_id, db_id, status, detail, error_message)
		VALUES (?, datetime('now'), ?, ?, ?, ?, ?, ?)`

	// eksekusi insert dengan semua parameter
	_, err := db.Exec(query,
		generateLogID(), // id unik untuk log ini
		event,           // jenis event, contoh: DB_CREATED
		userID,          // id user yang memicu event
		dbID,            // id database yang terlibat
		status,          // success atau failed
		detail,          // detail tambahan, contoh: port dan package
		errMsg,          // pesan error jika ada, kosong jika sukses
	)
	return err
}
