# SQL migrations

Versioned `*.sql` files applied in lexicographic order by the Go API at startup (or via `cmd/migrate`).

Set `SKIP_SQL_MIGRATIONS=1` to skip SQL migration step at API boot.
