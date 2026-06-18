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

## Railway / production

- The server listens on `HTTP_ADDR` if set, otherwise **`PORT`** (injected by Railway), else `:8080`.
- Do **not** hardcode `HTTP_ADDR=:8080` in production if the platform sets `PORT`.
- `CORS_ALLOWED_ORIGINS` must list exact origins **without quotes**, e.g. `https://management.woragis.me,https://www.woragis.me`.

## Docker (full stack)

```bash
docker compose up -d --build
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for domain layout and API conventions.

## API routes

| Route | Auth | Description |
|-------|------|-------------|
| `GET /health` | — | Liveness |
| `GET /ready` | — | Readiness (DB ping) |
| `GET /v1/admin/projects` | `X-Admin-Key` | List all projects |
| `POST /v1/admin/projects` | `X-Admin-Key` | Create project |
| `GET /v1/admin/media` | `X-Admin-Key` | List media assets |
| `POST /v1/admin/media` | `X-Admin-Key` | Upload file (multipart `file`) |
| `GET /v1/admin/profile` | `X-Admin-Key` | Get profile |
| `PATCH /v1/admin/profile` | `X-Admin-Key` | Update profile |
| `GET /v1/public/projects` | — | Public projects (`?featured=true`) |
| `GET /v1/public/projects/{slug}` | — | Public project detail |
| `GET /v1/public/profile` | — | Public profile |
| `GET /v1/public/media/{id}/file` | — | Serve media file |
