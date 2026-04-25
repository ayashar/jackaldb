# JackalDB API Service

REST API service untuk JackalDB yang dibangun dengan Node.js dan Express. Service ini berfungsi sebagai inti provisioning yang secara otomatis membuat dan menghapus instance PostgreSQL via Docker.

## Prerequisites

- Node.js 18+
- pnpm
  - Install: `npm install -g pnpm`
- Docker Desktop (harus dalam keadaan **running**)

## Setup

```bash
cd api
pnpm install
cp .env.example .env
```

## Menjalankan Service

```bash
# Mode production
pnpm start

# Mode development (auto-restart saat file berubah)
pnpm dev

# Server berjalan di http://localhost:3000
```

## Endpoints

- `POST /create-db` — Provision instance PostgreSQL baru
- `DELETE /delete-db/:id` — Hapus instance database
- `GET /databases` — List semua database yang sedang running
- `GET /databases/:id/status` — Cek status satu database
- `GET /health` — Cek apakah server sedang berjalan

### POST /create-db — Contoh request

```bash
curl -X POST http://localhost:3000/create-db \
  -H "Content-Type: application/json" \
  -d '{
    "db_name": "myapp",
    "package": "small",
    "user_id": "user_123"
  }'
```

Contoh response:

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

### DELETE /delete-db/:id — Contoh request

```bash
curl -X DELETE http://localhost:3000/delete-db/jkl_a3f9c2
```

### GET /databases — Contoh request

```bash
curl http://localhost:3000/databases
```

### GET /databases/:id/status — Contoh request

```bash
curl http://localhost:3000/databases/jkl_a3f9c2/status
```

## Paket Tersedia

| Paket | vCPU | RAM |
|-------|------|-----|
| `small` | 0.5 core | 256 MB |
| `medium` | 1 core | 512 MB |
| `large` | 2 core | 1 GB |

## Environment Variables

| Variable | Default | Keterangan |
|----------|---------|------------|
| `API_PORT` | `3000` | Port server |
| `DOCKER_SOCKET` | `/var/run/docker.sock` | Path ke Docker daemon |
| `PG_IMAGE` | `postgres:15` | Docker image PostgreSQL |
| `LOGGER_URL` | `http://localhost:8081/log` | URL logger service |

## Event Types

- `DB_CREATED` — Setelah container PostgreSQL berhasil di-provision
- `DB_DELETED` — Setelah container berhasil dihapus
- `PROVISION_FAILED` — Container gagal di-spawn

---

## Catatan untuk Pengguna Windows (PowerShell)

Ada beberapa penyesuaian yang perlu dilakukan jika menjalankan project ini di Windows menggunakan PowerShell.

### 1. Ganti Docker Socket di `.env`

Docker Desktop di Windows tidak menggunakan path yang sama seperti Mac/Linux. Buka file `.env` dan ganti nilai `DOCKER_SOCKET`:

```env
# Mac/Linux (default)
DOCKER_SOCKET=/var/run/docker.sock

# Windows — ganti menjadi ini
DOCKER_SOCKET=//./pipe/docker_engine
```

### 2. Pull image PostgreSQL terlebih dahulu

Docker tidak otomatis download image saat pertama kali dipakai. Jalankan perintah ini sekali sebelum menjalankan API:

```powershell
docker pull postgres:15
```

Tunggu hingga proses download selesai (~100-200 MB).

### 3. Format curl di PowerShell berbeda

PowerShell tidak mendukung single quote `'` untuk JSON. Gunakan `curl.exe` dengan double quote dan escape menggunakan backtick `` ` ``:

```powershell
curl.exe -X POST http://localhost:3000/create-db -H "Content-Type: application/json" -d "{`"db_name`":`"test`",`"package`":`"small`",`"user_id`":`"user_123`"}"
```

```powershell
curl.exe -X DELETE http://localhost:3000/delete-db/jkl_a3f9c2
```

```powershell
curl.exe http://localhost:3000/databases
```

> **Catatan:** Gunakan `curl.exe` bukan `curl` di PowerShell — `curl` di PowerShell adalah alias untuk `Invoke-WebRequest` yang berbeda perilakunya.