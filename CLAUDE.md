# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Start Postgres (Docker required)
make postgres

# Run database migrations
make migrateup       # apply all migrations
make migratedown     # roll back all migrations
make migrateup1      # apply one migration
make migratedown1    # roll back one migration

# Run the server
make server          # go run main.go

# Run all tests
make test            # go test -count=1 ./...

# Run a single test package
go test -count=1 ./api/...
go test -count=1 ./database/sqlc/...

# Run a single test
go test -count=1 -run TestCreateAccount ./api/...

# Regenerate sqlc code after editing SQL queries
make sqlc

# Regenerate mock after editing database/sqlc/store.go
make mockdb
```

## Architecture

**Moneytor** is a personal finance tracker. The backend is a Go REST API backed by PostgreSQL; the frontend is plain HTML/CSS/JS served separately via VS Code Live Server (expects origin `http://localhost:5500`).

### Layer overview

```
main.go               — wires config, DB pool, Store, and Server
utils/config.go       — loads config.json (Env, DBSource, HttpServerAddress)
api/                  — Gin HTTP handlers; one file per resource
database/
  migrations/         — golang-migrate SQL files (up/down)
  queries/            — raw SQL used by sqlc to generate Go code
  sqlc/               — generated code (DO NOT edit *.sql.go)
    store.go          — Store interface + SQLStore; custom transactions go here
    tx_create_entry.go — example transaction: creates entry + updates balance
  mocks/store.go      — mockgen-generated mock of Store (used in api tests)
frontend/             — static HTML/CSS/JS for the accounts view
```

### Key design patterns

- **Store interface** (`database/sqlc/store.go`): `SQLStore` embeds `*Queries` (sqlc-generated) and adds transaction methods. All API handlers depend on the `Store` interface, enabling mock injection in tests.
- **Mock-based API tests**: `api/*_test.go` files use `mockdb.NewMockStore` (from `database/mocks/`) rather than a real DB. DB-layer tests (`database/sqlc/*_test.go`) hit a real Postgres instance configured via `config.json`.
- **Amount sign convention**: expense entries are stored with a negative `amount`; income entries with positive. `ResolverEntryAmount` in `api/entries.go` applies this using the `TransactionType` constants (Expense=1, Income=2, Transfer=3).
- **Monetary amounts as integers**: `amount` is stored as `int64` (smallest currency unit, e.g. cents). No `decimal` type is used in current models despite the sqlc override in `sqlc.yaml`.
- **CORS**: hardcoded to allow `http://localhost:5500` and `http://127.0.0.1:5500` (VS Code Live Server). Update `api/server.go` `allowedOrigins` for other dev setups.

### Adding a new resource

1. Write SQL in `database/queries/<resource>.sql`
2. Run `make sqlc` to regenerate `database/sqlc/`
3. Add handler file `api/<resource>.go` with request structs and Gin handlers
4. Register routes in `api/server.go`
5. If tests need DB interaction, add `database/sqlc/<resource>_test.go`; if API-layer only, add `api/<resource>_test.go` with mock expectations
6. Run `make mockdb` if `Store` interface changed
