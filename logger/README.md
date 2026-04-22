# JackalDB Logger Service

Activity logging service untuk JackalDB yang dibangun dengan Go dan SQLite. Service ini berfungsi sebagai audit trail yang mencatat semua event yang terjadi di sistem.

## Prerequisites

- Go 1.20+
- GCC 
  - Mac: `xcode-select --install`
  - Linux: `sudo apt-get install gcc`

## Setup

```bash
cd logger
go mod download
```

## Menjalankan Service

### Mode HTTP Server
```bash
go run main.go logger.go database.go
# Server berjalan di http://localhost:8081
```

### Mode CLI
```bash
# Lihat semua log
go run main.go logger.go database.go --list

# Filter by user
go run main.go logger.go database.go --list --user user_123

# Filter by event
go run main.go logger.go database.go --list --event DB_CREATED

# Limit jumlah hasil
go run main.go logger.go database.go --list --limit 10
```

## Endpoints
- POST /log: Terima event dari API service 
- GET /logs: Ambil semua log (support filter) 

### POST /log — Contoh request
```bash
curl -X POST http://localhost:8081/log \
  -H "Content-Type: application/json" \
  -d '{
    "Event": "DB_CREATED",
    "UserID": "user_123",
    "DbID": "jkl_a3f9c2",
    "Status": "success",
    "Detail": "Container spawned on port 54321, image: postgres:15, package: medium"
  }'
```

### GET /logs — Contoh request
```bash
# Semua log
curl "http://localhost:8081/logs"

# Filter by user
curl "http://localhost:8081/logs?user=user_123"

# Filter by event
curl "http://localhost:8081/logs?event=DB_CREATED"
```

## Event Types
- DB_CREATED: Setelah container PostgreSQL berhasil di-provision
- DB_DELETED Setelah container berhasil dihapus
- DB_ACCESSED User query status database
- PROVISION_FAILED Container gagal di-spawn 
