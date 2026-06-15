# Architecture — Woragis Management Backend

Go API following the Lingo pattern: handler → service → repository.

## Domains

| Domain | Purpose |
|--------|---------|
| `devproject` | Programming projects, links, domains, secrets |
| `media` | Image/file upload and public URLs |
| `profile` | Landing page hero / about (singleton) |

## API surfaces

| Prefix | Auth | Purpose |
|--------|------|---------|
| `/health`, `/ready` | none | Probes |
| `/v1/admin/*` | `X-Admin-Key` | Full CRUD (management frontend) |
| `/v1/public/*` | none | Read-only data for landing page |

## Layout

```text
backend/
├── migrations/
├── server/
│   ├── cmd/server/main.go
│   ├── cmd/migrate/main.go
│   └── internal/
│       ├── httpserver/
│       ├── middleware/
│       ├── apperrors/
│       ├── models/
│       ├── migrate/
│       ├── platform/postgres/
│       ├── devproject/
│       ├── media/
│       └── profile/
└── docker-compose.yml
```

## Conventions

- UUID primary keys, JSON camelCase, Postgres snake_case
- Errors: `{ "code", "message" }` via `apperrors`
- SQL migrations in `migrations/` + GORM AutoMigrate for domain tables
