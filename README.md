# JackalDB — Instant PostgreSQL as a Service

## Prerequisites
- Node.js 18+ (untuk API)
- Go 1.20+ (untuk Logger)
- Docker & Docker Daemon running
- SQLite3

## Setup

### 1. API Service (Node.js)
\`\`\`bash
cd api/
npm install
cp .env.example .env
npm run dev
# Server running di http://localhost:3000
\`\`\`

### 2. Logger Service (Golang)
\`\`\`bash
cd logger/
go mod download
go run main.go
# Server running di http://localhost:8081
# CLI: go run main.go --list
\`\`\`

## Testing

### Create Database
\`\`\`bash
curl -X POST http://localhost:3000/create-db \
  -H "Content-Type: application/json" \
  -d '{"db_name": "test_db", "package": "small", "user_id": "user_123"}'
\`\`\`

### List Logs (CLI)
\`\`\`bash
cd logger/
go run main.go --list --user user_123
\`\`\`

### List Logs (API)
\`\`\`bash
curl http://localhost:3000/logs?user=user_123
\`\`\`
