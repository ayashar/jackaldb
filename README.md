# JackalDB — Instant PostgreSQL, Zero Config

> Mini Project — Cloud Computing (MII212610)
> Semester Genap 2025/2026
> Departemen Ilmu Komputer dan Elektronika, FMIPA
> Universitas Gadjah Mada

---

## Tentang Proyek

JackalDB adalah layanan cloud berbasis **PaaS (Platform as a Service)** yang memungkinkan pengguna mendapatkan instance database PostgreSQL terisolasi secara instan hanya melalui satu API call. Pengguna tidak perlu melakukan instalasi server, konfigurasi database, atau manajemen infrastruktur secara manual.

Seluruh proses provisioning dijalankan otomatis melalui pendekatan **Operations as Code (OaC)** — setiap operasi infrastruktur (membuat container, mengalokasikan resource, generate credentials) direpresentasikan sebagai kode yang dapat dieksekusi secara programatik.

Proyek ini dibuat sebagai implementasi nyata konsep cloud computing, meniru cara kerja layanan seperti Supabase, Railway, dan Neon dalam skala kecil.

---

## Tim Pengembang — Kelompok E

| Nama | NIM | Role |
|------|-----|------|
| Ayasha Rahmadinni | 24/545462/PA/23178 | Engineer 1 — Node.js API & Dockerode |
| Aliya Khairun Nisa | 24/543832/PA/23111 | Engineer 2 — Golang Logger Service |
| Maulana Faris Al Ghifari | 24/544029/PA/23119 | Tech Lead & Dokumentasi |
| Satya Wira Pramudita | 24/543649/PA/23102 | QA & Testing |
| Widad Muhammad Rafi | 24/545635/PA/23190 | Dokumentasi |

**Dosen Pengampu:** I Gede Mujiyatna, S.Kom., M.Kom.

---

## Arsitektur Sistem

```
Client (curl / Postman)
        │
        ▼
┌───────────────────┐
│   REST API        │  Node.js + Express        :3000
│   (Engineer 1)    │  Dockerode wrapper
└────────┬──────────┘
         │ spawn / stop container
         ▼
┌───────────────────┐
│   Docker Engine   │  PostgreSQL containers
│                   │  port 54321, 54322, ...
└───────────────────┘
         │ POST /log (fire-and-forget)
         ▼
┌───────────────────┐
│   Logger Service  │  Golang + SQLite          :8081
│   (Engineer 2)    │  audit trail & CLI
└───────────────────┘
```

---

## Konsep Cloud yang Diimplementasikan

| Konsep | Implementasi |
|--------|-------------|
| PaaS | Pengguna hanya interaksi via API, infrastruktur tersembunyi |
| On-demand provisioning | Container dibuat hanya saat diminta |
| Resource isolation | Tiap database dapat CPU, RAM, dan port sendiri |
| Operations as Code | Semua operasi Docker dijalankan via kode, bukan manual |

---

## Paket Layanan

| Paket | vCPU | RAM | Keterangan |
|-------|------|-----|------------|
| `small` | 0.5 core | 256 MB | Untuk development & prototype |
| `medium` | 1 core | 512 MB | Untuk staging & testing |
| `large` | 2 core | 1 GB | Untuk production skala kecil |

---

## Struktur Repository

```
jackaldb/
├── api/                        # Engineer 1 — Node.js REST API
│   ├── src/
│   │   ├── app.js              # Entry point, Express setup & routes
│   │   ├── db-controller.js    # Handler tiap endpoint
│   │   ├── docker-service.js   # Dockerode wrapper, OaC logic
│   │   ├── credentials.js      # Generate username & password aman
│   │   └── error-handler.js    # Error handling middleware
│   ├── .env.example
│   ├── package.json
│   └── README.md
│
└── logger/                     # Engineer 2 — Golang Logger Service
    ├── main.go                 # HTTP server & CLI entry point
    ├── database.go             # SQLite schema & query
    ├── logger.go               # Insert log event
    ├── logs.db                 # SQLite database file
    ├── go.mod
    └── go.sum
```

---

## Prasyarat

Pastikan semua sudah terinstall sebelum menjalankan proyek:

- [Node.js](https://nodejs.org/) v18 atau lebih baru
- [pnpm](https://pnpm.io/) (`npm install -g pnpm`)
- [Go](https://go.dev/) v1.20 atau lebih baru
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (harus dalam keadaan **running**)
- [psql](https://www.postgresql.org/) — opsional, untuk verifikasi koneksi database

---

## Instalasi & Menjalankan

### 1. Clone Repository

```bash
git clone https://github.com/ayashar/jackaldb.git
cd jackaldb
```

### 2. Jalankan Logger Service (Engineer 2)

```bash
cd logger
go run .
```

Logger akan berjalan di `http://localhost:8081`. Jalankan ini **sebelum** API.

### 3. Setup & Jalankan API (Engineer 1)

Buka terminal baru:

```bash
cd api
pnpm install
cp .env.example .env
pnpm start
```

API akan berjalan di `http://localhost:3000`.

---

## API Reference

### `POST /create-db`

Provision instance PostgreSQL baru.

**Request Body:**
```json
{
  "db_name": "myapp",
  "package": "small",
  "user_id": "user_123"
}
```

**Response `201 Created`:**
```json
{
  "success": true,
  "data": {
    "db_id": "jkl_a3f9c2",
    "db_name": "myapp",
    "host": "localhost",
    "port": 54321,
    "username": "jkl_user_x7k9ab2c",
    "password": "147c00a2e6c88a...",
    "connection_string": "postgresql://jkl_user_x7k9ab2c:147c00...@localhost:54321/myapp",
    "created_at": "2026-04-25T07:02:03.925Z"
  }
}
```

---

### `DELETE /delete-db/:id`

Hentikan dan hapus container database.

**Contoh:**
```bash
curl -X DELETE http://localhost:3000/delete-db/jkl_a3f9c2
```

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "Database jkl_a3f9c2 berhasil dihapus"
}
```

---

### `GET /databases`

List semua container PostgreSQL yang sedang running.

```bash
curl http://localhost:3000/databases
```

---

### `GET /databases/:id/status`

Cek status satu database berdasarkan `db_id`.

```bash
curl http://localhost:3000/databases/jkl_a3f9c2/status
```

---

### `GET /health`

Cek apakah API server sedang berjalan.

```bash
curl http://localhost:3000/health
```

---

## Logger Service Reference

### HTTP Endpoint

```bash
# Ambil semua log
curl http://localhost:8081/logs

# Filter berdasarkan user
curl "http://localhost:8081/logs?user=user_123"

# Filter berdasarkan event
curl "http://localhost:8081/logs?event=DB_CREATED"

# Kombinasi filter
curl "http://localhost:8081/logs?user=user_123&event=DB_CREATED&limit=10"
```

### CLI

```bash
cd logger

# Tampilkan semua log
go run . --list

# Filter berdasarkan user
go run . --list --user user_123

# Filter berdasarkan event
go run . --list --event DB_CREATED

# Batasi jumlah hasil
go run . --list --limit 50
```

### Event yang Dicatat

| Event | Pemicu |
|-------|--------|
| `DB_CREATED` | Container berhasil di-spawn |
| `DB_DELETED` | Container berhasil dihapus |
| `PROVISION_FAILED` | Gagal membuat container |

---

## Environment Variables

Salin `.env.example` ke `.env` lalu sesuaikan:

```env
API_PORT=3000
DOCKER_SOCKET=/var/run/docker.sock
PG_IMAGE=postgres:15
LOGGER_URL=http://localhost:8081/log
```

---

## Testing End-to-End

```bash
# 1. Buat database
curl -X POST http://localhost:3000/create-db \
  -H "Content-Type: application/json" \
  -d '{"db_name":"test","package":"small","user_id":"user_123"}'

# 2. Konek ke database (ganti dengan connection_string dari response)
psql "postgresql://jkl_user_xxx:password@localhost:54321/test"

# 3. Jalankan query di dalam psql
SELECT 1;
SELECT current_database();

# 4. Cek log tercatat
cd logger && go run . --list

# 5. Hapus database
curl -X DELETE http://localhost:3000/delete-db/jkl_xxxxxx

# 6. Verifikasi container terhapus
docker ps | grep jackaldb
```

---

*Mini Project Cloud Computing — Kelompok E — Universitas Gadjah Mada 2026*