# Messaging: template catalog & data bindings

## Problem

Messaging had two parallel template systems (LeetCode `WhatsappMessageTemplate` vs generic `MessageTemplate`), jobs sent literal bodies without `{{var}}` resolution, skips were invisible, failed sends still advanced cron, and `ai_assisted` / Telegram were not wired.

## Model

```
Program (leetcode | project | custom)
  └── catalog fields (API: GET /v1/admin/messaging/catalog?program=…)

MessageTemplate
  ├── body (with {{placeholders}})
  ├── bindings: { "problemTitle": "leetcode.problemTitle", … }
  └── composeMode: static | ai_assisted

ScheduledJob
  ├── templateSlug + programAction
  └── dataSource: { program, date, projectId, projectSlug }
```

## Backend (T1–T3)

| Area | Change |
|------|--------|
| Deliveries | `sent` / `failed` / `skipped`; skip recorded; cron **not** advanced on `failed` |
| Executor | Template renderer, agent compose, Telegram worker |
| Content | `ResolveDispatchVars`; LeetCode dispatch prefers `MessageTemplate` when slug matches |
| API | `bindings` on templates, `dataSource` on jobs, catalog + preview endpoints |

## Resolvers

- **leetcode** — uses `programAction` (`problem`, `discussion`, `solution`, `weekly`) + optional `dataSource.date`
- **project** — requires `dataSource.projectId` or `projectSlug`

## Frontend

- Templates: bindings JSON, field catalog hint, preview
- Jobs: data source program + project id
- Deliveries: `skipped` status visible

## Env

- `AGENT_WORKER_URL` + `AGENT_API_KEY` — `ai_assisted` compose
- `TELEGRAM_WORKER_URL` + `WORKER_API_KEY` — Telegram send
- `WHATSAPP_WORKER_URL` — unchanged

## Follow-ups

- WhatsApp group sync from worker
- Retry/backoff policy for failed deliveries
- Migrate LeetCode templates fully into Messaging UI
