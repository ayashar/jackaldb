package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// logeentry adalah struktur data untuk satu baris log di database
type LogEntry struct {
	LogID     string `json:"LogID"`
	Timestamp string `json:"Timestamp"`
	Event     string `json:"Event"`
	UserID    string `json:"UserID"`
	DbID      string `json:"DbID"`
	Status    string `json:"Status"`
	Detail    string `json:"Detail"`
	ErrorMsg  string `json:"ErrorMsg"`
}

// initdatabase membuat tabel logs jika belum ada
// dipanggil sekali saat aplikasi pertama kali jalan
func initDatabase(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS logs (
			log_id        TEXT PRIMARY KEY,
			timestamp     DATETIME DEFAULT CURRENT_TIMESTAMP,
			event         TEXT NOT NULL,
			user_id       TEXT NOT NULL,
			db_id         TEXT,
			status        TEXT,
			detail        TEXT,
			error_message TEXT
		);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("gagal membuat tabel logs:", err)
	}
}

// getlogs mengambil data log dari sqlite dengan filter opsional
// bisa difilter berdasarkan user_id, event type, dan jumlah limit
func getLogs(db *sql.DB, userID, eventType string, limit int) []LogEntry {
	// query dasar, WHERE 1=1 supaya mudah tambah kondisi dinamis
	query := `SELECT log_id, timestamp, event, user_id,
		COALESCE(db_id,''), COALESCE(status,''), COALESCE(detail,''), COALESCE(error_message,'')
		FROM logs WHERE 1=1`

	var args []interface{}

	// tambah filter user_id jika diberikan
	if userID != "" {
		query += " AND user_id = ?"
		args = append(args, userID)
	}

	// tambah filter event jika diberikan
	if eventType != "" {
		query += " AND event = ?"
		args = append(args, eventType)
	}

	// urutkan dari terbaru, batasi jumlah hasil
	query += fmt.Sprintf(" ORDER BY timestamp DESC LIMIT %d", limit)

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal("error saat query logs:", err)
	}
	defer rows.Close() // pastikan rows ditutup setelah selesai

	var entries []LogEntry
	for rows.Next() {
		var e LogEntry
		// scan tiap kolom ke field struct
		if err := rows.Scan(
			&e.LogID, &e.Timestamp, &e.Event, &e.UserID,
			&e.DbID, &e.Status, &e.Detail, &e.ErrorMsg,
		); err != nil {
			log.Println("error saat scan baris log:", err)
			continue // lewati baris yang error, jangan hentikan semua
		}
		entries = append(entries, e)
	}
	return entries
}
