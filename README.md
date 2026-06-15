# Woragis Management Backend

Personal management API for programming projects, media assets, and profile data consumed by the landing page.

## Stack

- Go 1.24, stdlib HTTP, GORM, PostgreSQL
- Pattern: handler → service → repository (Lingo-style)

## Local development

```bash
# Start Postgres
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d postgres

# Copy env
cp .env.example .env

# Run API (from server/)
cd server
go run ./cmd/server
```

Health: `GET http://127.0.0.1:8080/health`

## Docker (full stack)

```bash
docker compose up -d --build
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for domain layout and API conventions.
