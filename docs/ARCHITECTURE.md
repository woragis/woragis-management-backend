# Architecture — Woragis Management Backend

Go API following the Lingo pattern: handler → service → repository.

**Planos detalhados (agente, contatos, finanças, personalidade):** [docs/plans/README.md](./plans/README.md)

## Domains

| Domain | Purpose |
|--------|---------|
| `devproject` | Programming projects, links, domains, secrets |
| `media` | Image/file upload and public URLs |
| `profile` | Landing page hero / about (singleton) |
| `finance` | Income, expenses, transactions, invoices, budgets |
| `content` | LeetCode videos, thumbnails, WhatsApp templates |

### Planned

| Domain | Doc |
|--------|-----|
| `contacts` | [plans/02-contacts-domain.md](./plans/02-contacts-domain.md) |
| `agent` (personality, tools) | [plans/04-agent-personality.md](./plans/04-agent-personality.md) |

## API surfaces

| Prefix | Auth | Purpose |
|--------|------|---------|
| `/health`, `/ready` | none | Probes |
| `/v1/admin/*` | `X-Admin-Key` | Full CRUD (management frontend) |
| `/v1/public/*` | none | Read-only data for landing page |
| `/v1/internal/content/leetcode/*` | `X-Worker-Key` | WhatsApp worker (LeetCode dispatch) |
| `/v1/internal/agent/*` | `AGENT_API_KEY` | Planned — agent tools |
| `/v1/webhooks/creatives` | none | Creatives thumbnail callback |

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
│       ├── profile/
│       ├── finance/
│       └── content/
├── agent-worker/     # planned — see docs/plans/
├── whatsapp-worker/
└── docker-compose.yml
```

## Conventions

- UUID primary keys, JSON camelCase, Postgres snake_case
- Errors: `{ "code", "message" }` via `apperrors`
- SQL migrations in `migrations/` + GORM AutoMigrate for domain tables
