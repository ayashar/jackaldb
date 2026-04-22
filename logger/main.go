package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// db adalah koneksi global ke sqlite, dipakai oleh semua handler
var db *sql.DB

func main() {
	// definisi flag untuk mode cli
	listFlag := flag.Bool("list", false, "tampilkan semua log")
	userFlag := flag.String("user", "", "filter log berdasarkan user id")
	eventFlag := flag.String("event", "", "filter log berdasarkan jenis event")
	limitFlag := flag.Int("limit", 20, "batasi jumlah log yang ditampilkan")
	flag.Parse()

	// buka koneksi ke file sqlite
	var err error
	db, err = sql.Open("sqlite3", "./logs.db")
	if err != nil {
		log.Fatal("tidak bisa membuka database:", err)
	}
	defer db.Close()

	// inisialisasi tabel jika belum ada
	initDatabase(db)

	// mode cli: tampilkan log lalu keluar
	if *listFlag {
		logs := getLogs(db, *userFlag, *eventFlag, *limitFlag)
		if len(logs) == 0 {
			fmt.Println("tidak ada log ditemukan.")
			return
		}
		// cetak tiap log dalam format yang mudah dibaca
		for _, l := range logs {
			fmt.Printf("[%s] %s | user=%s | db=%s | status=%s | %s\n",
				l.Timestamp, l.Event, l.UserID, l.DbID, l.Status, l.Detail)
		}
		return
	}

	// mode http server: daftarkan endpoint dan mulai listen
	http.HandleFunc("/log", handleIncomingLog) // menerima event dari node.js
	http.HandleFunc("/logs", handleGetLogs)    // melayani query log
	fmt.Println("logger service berjalan di :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// handleincominglog menangani POST /log dari engineer 1 (node.js)
// setiap kali create-db atau delete-db berhasil, node.js tembak ke sini
func handleIncomingLog(w http.ResponseWriter, r *http.Request) {
	// hanya terima method POST
	if r.Method != http.MethodPost {
		http.Error(w, "method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	// decode body json ke struct LogEntry
	var entry LogEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "format json tidak valid", http.StatusBadRequest)
		return
	}

	// simpan event ke sqlite
	if err := LogEvent(db, entry.Event, entry.UserID, entry.DbID,
		entry.Status, entry.Detail, entry.ErrorMsg); err != nil {
		http.Error(w, "gagal menyimpan log: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// kirim response sukses ke node.js
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "logged"})
}

// handlegetlogs menangani GET /logs dengan filter opsional
// bisa diakses via rest api atau dari frontend untuk monitoring
func handleGetLogs(w http.ResponseWriter, r *http.Request) {
	// ambil parameter filter dari query string
	userID := r.URL.Query().Get("user")
	eventType := r.URL.Query().Get("event")

	// default limit 20, bisa dioverride via query param
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	// ambil data dari sqlite
	logs := getLogs(db, userID, eventType, limit)

	// kirim response dalam format json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"total":   len(logs),
		"data":    logs,
	})
}
