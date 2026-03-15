# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Moneytor is a personal finance web application with a Go REST API backend (Gin) and a vanilla JS frontend. PostgreSQL is the database, accessed via sqlc-generated type-safe code.

## Common Commands

All key commands are in the `Makefile`:

```bash
make server        # Start the API server (go run main.go)
make test          # Run all Go tests (go test -count=1 ./...)
make postgres      # Start PostgreSQL 17 in Docker
make migrateup     # Apply all pending migrations
make migratedown   # Roll back all migrations
make migrateup1    # Apply one migration
make migratedown1  # Roll back one migration
make sqlc          # Regenerate Go code from SQL queries
make mockdb        # Regenerate mocks from the Store interface
make migration     # Create a new migration file
```

## Architecture

### Backend (Go)

- **Entry point**: `main.go` — loads `config.json`, creates pgxpool, starts server
- **HTTP layer**: `api/` — Gin handlers; `api/server.go` sets up router, CORS middleware, and routes under `/api/v1`
- **Database layer**: `database/sqlc/` — sqlc-generated code; `store.go` defines the `Store` interface and a custom `execTx` helper for transactions
- **Queries**: SQL definitions live in `database/queries/`; run `make sqlc` after editing them to regenerate Go code
- **Mocks**: `database/mocks/store.go` is generated from the `Store` interface via `make mockdb`

### Database

- Migrations are in `database/migrations/` (managed with golang-migrate)
- sqlc is configured with `emit_json_tags: true` and `json_tags_case_style: "camel"` — all generated structs use camelCase JSON tags
- Monetary amounts are stored as `int64` (cents or smallest currency unit)
- Transaction entries reference a single `from_account_id`; `to_account_id` is nullable (used for transfers)
- `CreateEntryTx` in `store.go` is the main transaction: inserts the entry and updates the account balance atomically

### API

- All routes are under `/api/v1`
- CORS allows `http://localhost:5500` and `http://127.0.0.1:5500` (Live Server default ports)
- Pagination uses `pageId` + `pageSize` query params (default 5, max 10)
- PostgreSQL FK violations return HTTP 403; invalid input returns HTTP 400

### Frontend

- Currently vanilla HTML/CSS/JS in `frontend/`
- Fetches from `http://localhost:8080/api/v1`
- No build step required — open directly with a static file server (e.g. VS Code Live Server on port 5500)

## Configuration

`config.json` (gitignored) is required at runtime:

```json
{
  "Env": "DEV",
  "DBSource": "postgres://root:pass.123@localhost:5432/moneytor?sslmode=disable",
  "HttpServerAddress": "0.0.0.0:8080"
}
```

`Env: "DEV"` enables human-readable zerolog console output.
